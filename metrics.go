package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Metrics struct {
	mu           sync.RWMutex
	responseTime map[string][]time.Duration
	errors       map[string]int
	throughputs  map[string]int
}

func InitializeMetrics() *Metrics {
	return &Metrics{
		responseTime: make(map[string][]time.Duration),
		errors:       make(map[string]int),
		throughputs:  make(map[string]int),
	}
}

func (m *Metrics) Collect(path string, responseTime time.Duration, status int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.responseTime[path] = append(m.responseTime[path], responseTime)

	if status != http.StatusOK {
		m.errors[path]++
	}

	m.throughputs[path]++
}

func (m *Metrics) LogMetrics(interval time.Duration) {
	go func() {
		for {
			<-time.After(interval)

			m.mu.RLock()
			defer m.mu.RUnlock()

			for path, times := range m.responseTime {
				var totalDuration time.Duration
				for _, t := range times {
					totalDuration += t
				}
				avgDuration := totalDuration / time.Duration(len(times))
				log.Printf("Path: %s - Avg Response Time: %v, Errors: %d, Throughputs: %d\n", path, avgDuration, m.errors[path], m.throughputs[path])
			}
		}
	}()
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &wrappedWriter{ResponseWriter: w}
		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		m.Collect(r.URL.Path, duration, ww.status)
	})
}

type wrappedWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *wrappedWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.status = code
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		log.Fatal("API_PORT is not set")
	}

	metrics := InitializeMetrics()
	// Log metrics every 5 minutes
	metrics.LogMetrics(5 * time.Minute)

	http.Handle("/", metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})))

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}