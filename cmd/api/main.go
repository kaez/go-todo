package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kaez/go-todo/internal/handlers"
	"github.com/kaez/go-todo/internal/middleware"
	"github.com/kaez/go-todo/internal/repository"
)

func main() {
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "./todos.db")

	log.Printf("Starting Todo API on port %s", port)
	log.Printf("Database path: %s", dbPath)

	repo, err := repository.NewTodoRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close()

	todoHandler := handlers.NewTodoHandler(repo)
	healthHandler := handlers.NewHealthHandler(repo)
	metricsHandler := handlers.NewMetricsHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/todos", todoHandler.GetAll)
	mux.HandleFunc("GET /api/todos/{id}", todoHandler.GetByID)
	mux.HandleFunc("POST /api/todos", todoHandler.Create)
	mux.HandleFunc("PUT /api/todos/{id}", todoHandler.Update)
	mux.HandleFunc("DELETE /api/todos/{id}", todoHandler.Delete)

	mux.HandleFunc("GET /health/live", healthHandler.Liveness)
	mux.HandleFunc("GET /health/ready", healthHandler.Readiness)

	mux.HandleFunc("GET /metrics", metricsHandler.ServeMetrics)

	handler := middleware.Logging(mux)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := server.Close(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
