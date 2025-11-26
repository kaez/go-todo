package handlers

import (
	"fmt"
	"net/http"
	"runtime"
	"sync/atomic"
)

type MetricsHandler struct {
	requestCount atomic.Int64
	todoCount    atomic.Int64
}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

func (m *MetricsHandler) IncrementRequests() {
	m.requestCount.Add(1)
}

func (m *MetricsHandler) SetTodoCount(count int64) {
	m.todoCount.Store(count)
}

func (m *MetricsHandler) ServeMetrics(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	fmt.Fprintf(w, "# HELP http_requests_total The total number of HTTP requests\n")
	fmt.Fprintf(w, "# TYPE http_requests_total counter\n")
	fmt.Fprintf(w, "http_requests_total %d\n\n", m.requestCount.Load())

	fmt.Fprintf(w, "# HELP todos_total The total number of todos\n")
	fmt.Fprintf(w, "# TYPE todos_total gauge\n")
	fmt.Fprintf(w, "todos_total %d\n\n", m.todoCount.Load())

	fmt.Fprintf(w, "# HELP go_memstats_alloc_bytes Number of bytes allocated\n")
	fmt.Fprintf(w, "# TYPE go_memstats_alloc_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_alloc_bytes %d\n\n", memStats.Alloc)

	fmt.Fprintf(w, "# HELP go_goroutines Number of goroutines\n")
	fmt.Fprintf(w, "# TYPE go_goroutines gauge\n")
	fmt.Fprintf(w, "go_goroutines %d\n", runtime.NumGoroutine())
}
