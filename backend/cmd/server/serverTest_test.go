package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestHealthAndReadiness(t *testing.T) {
	handler := newHandler(slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, path := range []string{"/healthz", "/readyz"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "status") {
			t.Fatalf("unexpected response for %s", path)
		}
	}
}

func TestWebSocketConnectionEnvelope(t *testing.T) {
	server := httptest.NewServer(newHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	connection, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseNow()
	_, payload, err := connection.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(payload), "connection.accepted") || !strings.Contains(string(payload), "1.0.0") {
		t.Fatalf("unexpected connection event: %s", payload)
	}
}
