package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	matchdomain "crownfall/backend/internal/game/domain/match"
	platformrandom "crownfall/backend/internal/platform/random"
	"crownfall/backend/internal/realtime/rooms"
	"crownfall/backend/internal/realtime/websocket"
)

func newHandler(logger *slog.Logger, registry *rooms.Registry) http.Handler {
	mux := http.NewServeMux()
	realtime := websocket.NewServer(registry)
	response := func(status string) http.HandlerFunc {
		return func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(map[string]string{"status": status})
		}
	}
	mux.HandleFunc("GET /healthz", response("ok"))
	mux.HandleFunc("GET /readyz", response("ready"))
	mux.HandleFunc("POST /v1/rooms", func(writer http.ResponseWriter, request *http.Request) {
		seed, err := (platformrandom.CryptoSource{}).Seed()
		if err != nil {
			http.Error(writer, "random source unavailable", http.StatusServiceUnavailable)
			return
		}
		roomID := "room-" + platformrandom.Commitment(seed)[:12]
		state := matchdomain.New(roomID, seed, platformrandom.Commitment(seed))
		if !registry.Add(roomID, rooms.NewRoom(state, 64)) {
			http.Error(writer, "room collision", http.StatusConflict)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(writer).Encode(map[string]string{"roomId": roomID})
	})
	mux.HandleFunc("GET /ws", realtime.Handler())
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		started := time.Now()
		mux.ServeHTTP(writer, request)
		logger.Info("http request", "method", request.Method, "path", request.URL.Path, "duration_ms", time.Since(started).Milliseconds())
	})
}
