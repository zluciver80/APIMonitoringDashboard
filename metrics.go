package main

import (
    "log"
    "net/http"
    "os"
    "sync"
    "time"
)

type APIMetrics struct {
    mutex        sync.RWMutex
    responseDurations map[string][]time.Duration
    errorCounts  map[string]int
    requestCounts map[string]int
}

func NewAPIMetrics() *APIMetrics {
    return &APIMetrics{
        responseDurations: make(map[string][]time.Duration),
        errorCounts:      make(map[string]int),
        requestCounts:    make(map[string]int),
    }
}

func (m *APIMetrics) RecordRequestMetrics(path string, duration time.Duration, statusCode int) {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    m.responseDurations[path] = append(m.responseDurations[path], duration)

    if statusCode != http.StatusOK {
        m.errorCounts[path]++
    }

    m.requestCounts[path]++
}

func (m *APIMetrics) PeriodicallyLogMetrics(interval time.Duration) {
    go func() {
        for {
            <-time.After(interval)

            m.mutex.RLock()
            func() {
                defer m.mutex.RUnlock()

                for path, durations := range m.responseDurations {
                    var totalDuration time.Duration
                    for _, d := range durations {
                        totalDuration += d
                    }
                    avgDuration := totalDuration / time.Duration(len(durations))
                    log.Printf("Path: %s - Avg Response Time: %v, Error Count: %d, Request Count: %d\n",
                        path, avgDuration, m.errorCounts[path], m.requestCounts[path])
                }
            }()
        }
    }()
}

func (m *APIMetrics) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        startTime := time.Now()

        wrappedResponseWriter := &responseWriterWrapper{ResponseWriter: w}
        next.ServeHTTP(wrappedResponseWriter, r)

        duration := time.Since(startTime)

        if !wrappedResponseWriter.headerWritten {
            wrappedResponseWriter.statusCode = http.StatusOK
        }

        m.RecordRequestMetrics(r.URL.Path, duration, wrappedResponseWriter.statusCode)
    })
}

type responseWriterWrapper struct {
    http.ResponseWriter
    statusCode   int
    headerWritten bool
}

func (w *responseWriterWrapper) WriteHeader(code int) {
    if !w.headerWritten {
        w.statusCode = code
        w.headerWritten = true
        w.ResponseWriter.WriteHeader(code)
    }
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
    if !w.headerWritten {
        w.WriteHeader(http.StatusOK)
    }
    return w.ResponseWriter.Write(b)
}

func main() {
    port := os.Getenv("API_PORT")
    if port == "" {
        log.Fatal("API_PORT is not set")
    }

    apiMetrics := NewAPIMetrics()
    apiMetrics.PeriodicallyLogMetrics(5 * time.Minute)

    http.Handle("/", apiMetrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, err := w.Write([]byte("Hello, world!"))
        if err != nil {
            log.Printf("Error writing response: %v", err)
        }
    })))

    log.Printf("Starting server on port %s...", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}