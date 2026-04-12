package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewHealthHandler(db *pgxpool.Pool, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, _ *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	checks := map[string]string{
		"postgres": "ok",
		"redis":    "ok",
	}

	status := http.StatusOK
	if err := h.db.Ping(ctx); err != nil {
		checks["postgres"] = "failed"
		status = http.StatusServiceUnavailable
	}
	if err := h.redis.Ping(ctx).Err(); err != nil {
		checks["redis"] = "failed"
		status = http.StatusServiceUnavailable
	}

	responseStatus := "ready"
	if status != http.StatusOK {
		responseStatus = "not_ready"
	}

	respondWithJSON(w, status, map[string]any{
		"status": responseStatus,
		"checks": checks,
	})
}
