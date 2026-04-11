package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rostmebel/backend/internal/config"
	"github.com/rostmebel/backend/internal/interfaces/http/handler"
	"github.com/rostmebel/backend/internal/interfaces/http/middleware"
)

type Server struct {
	router *chi.Mux
	port   string
}

func NewServer(cfg *config.Config, ph *handler.ProductHandler, oh *handler.OrderHandler, ah *handler.AdminHandler) *Server {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		// Public
		r.Get("/products", ph.GetProducts)
		r.Get("/products/{id}", ph.GetProduct)
		r.Get("/categories", ph.GetCategories)
		r.Post("/orders", oh.CreateOrder)
		r.Post("/ai/search", ph.AISearch)

		// Admin Auth
		r.Post("/admin/auth/login", ah.Login)
		r.Post("/admin/auth/refresh", ah.Refresh)

		// Protected Admin
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWTSecret))
			
			r.Post("/admin/auth/logout", ah.Logout)
			r.Get("/admin/stats", ah.GetStats)

			r.Route("/admin/products", func(r chi.Router) {
				r.Get("/", ph.GetProducts)
				r.Post("/", ph.CreateProduct)
				r.Get("/export", ph.ExportProducts)
				r.Put("/{id}", ph.UpdateProduct)
				r.Delete("/{id}", ph.DeleteProduct)
			})

			r.Post("/admin/upload", ph.UploadImage)

			r.Route("/admin/orders", func(r chi.Router) {
				r.Get("/", oh.GetOrders)
				r.Patch("/{id}/status", oh.UpdateOrderStatus)
				r.Post("/{id}/spam", oh.MarkAsSpam)
				r.Get("/export", oh.ExportOrders)
			})
		})
	})

	return &Server{
		router: r,
		port:   cfg.Port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	fmt.Printf("Server starting on port %s\n", s.port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
