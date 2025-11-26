package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/todo-api/internal/repository"
)

type HealthHandler struct {
	repo *repository.TodoRepository
}

func NewHealthHandler(repo *repository.TodoRepository) *HealthHandler {
	return &HealthHandler{repo: repo}
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"

	if _, err := h.repo.GetAll(); err != nil {
		dbStatus = "error"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:   "ready",
		Database: dbStatus,
	})
}
