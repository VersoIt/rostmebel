package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rostmebel/backend/internal/application/admin"
	"github.com/rostmebel/backend/internal/application/order"
	"github.com/rostmebel/backend/internal/application/product"
	"github.com/rostmebel/backend/internal/application/review"
	"github.com/rostmebel/backend/internal/config"
	domAdmin "github.com/rostmebel/backend/internal/domain/admin"
	"github.com/rostmebel/backend/internal/infrastructure/gemini"
	"github.com/rostmebel/backend/internal/infrastructure/httpx"
	"github.com/rostmebel/backend/internal/infrastructure/postgres"
	"github.com/rostmebel/backend/internal/infrastructure/redis"
	"github.com/rostmebel/backend/internal/infrastructure/telegram"
	"github.com/rostmebel/backend/internal/interfaces/http"
	"github.com/rostmebel/backend/internal/interfaces/http/handler"
	"github.com/rostmebel/backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)
	warnUnsafeProductionConfig(cfg, log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Infrastructure
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Error("failed to run migrations", "error", err)
		if err != migrate.ErrNoChange {
			os.Exit(1)
		}
	}

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	rdb, err := redis.NewClient(cfg.RedisURL, cfg.RedisPassword)
	if err != nil {
		log.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}

	outboundHTTPClient, err := httpx.NewHTTPClient(httpx.ClientOptions{
		Timeout:               30 * time.Second,
		DialTimeout:           10 * time.Second,
		KeepAlive:             30 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		MaxConnsPerHost:       50,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy: &httpx.ProxyConfig{
			Scheme:   cfg.OutboundProxyScheme,
			Host:     cfg.OutboundProxyHost,
			Port:     cfg.OutboundProxyPort,
			Username: cfg.OutboundProxyUsername,
			Password: cfg.OutboundProxyPassword,
		},
	})
	if err != nil {
		log.Error("failed to create outbound http client", "error", err)
		os.Exit(1)
	}

	geminiClient := gemini.NewClientWithOptions(gemini.ClientOptions{
		APIKey:         cfg.GeminiAPIKey,
		Model:          cfg.GeminiModel,
		FallbackModels: cfg.GeminiFallbackModels,
		HTTPClient:     outboundHTTPClient,
	})
	tgClient := telegram.NewClient(cfg.TelegramToken, cfg.TelegramChatID, outboundHTTPClient)

	// Repositories
	productRepo := postgres.NewProductRepo(pool)
	orderRepo := postgres.NewOrderRepo(pool)
	adminRepo := postgres.NewAdminRepo(pool)
	reviewRepo := postgres.NewReviewRepo(pool)

	// Seed first admin
	if err := seedAdmin(ctx, adminRepo, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		log.Error("failed to seed admin", "error", err)
		os.Exit(1)
	}

	// UseCases
	productUC := product.NewUseCase(productRepo)
	aiUC := product.NewAIUseCase(productRepo, geminiClient, rdb, log)
	orderUC := order.NewUseCase(orderRepo, productRepo, rdb, tgClient, cfg.OrderLimitEnabled)
	adminUC := admin.NewUseCase(adminRepo, cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	reviewUC := review.NewUseCase(reviewRepo, orderRepo)

	// Handlers
	ph := handler.NewProductHandler(productUC, aiUC, cfg.PublicSiteURL)
	oh := handler.NewOrderHandler(orderUC)
	ah := handler.NewAdminHandler(adminUC)
	rh := handler.NewReviewHandler(reviewUC)
	hh := handler.NewHealthHandler(pool, rdb)

	// Server
	srv := http.NewServer(cfg, ph, oh, ah, rh, hh)

	if err := srv.Start(ctx); err != nil {
		log.Error("server error", "error", err)
	}
}

func seedAdmin(ctx context.Context, repo domAdmin.Repository, username, password string) error {
	a, err := repo.GetByUsername(ctx, username)
	if err != nil || a != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newAdmin := &domAdmin.Admin{
		Username:     username,
		PasswordHash: string(hash),
	}
	return repo.Create(ctx, newAdmin)
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func warnUnsafeProductionConfig(cfg *config.Config, log *slog.Logger) {
	if cfg.AppEnv != "production" {
		return
	}
	if cfg.JWTSecret == "" || cfg.JWTSecret == "default-secret" || len(cfg.JWTSecret) < 32 {
		log.Warn("JWT_SECRET is weak for production; rotate it to a long random value")
	}
	if strings.EqualFold(cfg.AdminUsername, "admin") && cfg.AdminPassword == "admin" {
		log.Warn("default admin credentials are unsafe for production")
	}
}
