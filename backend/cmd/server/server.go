package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"crownfall/backend/internal/realtime/websocket"
)

func newHandler(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	response := func(status string) http.HandlerFunc {
		return func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(map[string]string{"status": status})
		}
	}
	mux.HandleFunc("GET /healthz", response("ok"))
	mux.HandleFunc("GET /readyz", response("ready"))
	mux.HandleFunc("GET /ws", websocket.Handler())
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		started := time.Now()
		mux.ServeHTTP(writer, request)
		logger.Info("http request", "method", request.Method, "path", request.URL.Path, "duration_ms", time.Since(started).Milliseconds())
	})
}
