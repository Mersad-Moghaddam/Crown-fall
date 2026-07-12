package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"crownfall/backend/internal/game/engine"
	"crownfall/backend/internal/realtime/rooms"
	"github.com/coder/websocket"
)

func TestHealthReadinessAndRoomCreation(t *testing.T) {
	registry := rooms.NewRegistry()
	defer registry.Close()
	handler := newHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), registry)
	for _, path := range []string{"/healthz", "/readyz"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("unexpected response for %s", path)
		}
	}
	request := httptest.NewRequest(http.MethodPost, "/v1/rooms", bytes.NewReader(nil))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create room status: %d", response.Code)
	}
	var payload map[string]string
	_ = json.Unmarshal(response.Body.Bytes(), &payload)
	if _, exists := registry.Get(payload["roomId"]); !exists {
		t.Fatal("created room was not registered")
	}
}

func TestSixClientWebSocketBootstrapAndReconnect(t *testing.T) {
	registry := rooms.NewRegistry()
	defer registry.Close()
	server := httptest.NewServer(newHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), registry))
	defer server.Close()
	response, err := http.Post(server.URL+"/v1/rooms", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	var created map[string]string
	_ = json.NewDecoder(response.Body).Decode(&created)
	response.Body.Close()
	roomID := created["roomId"]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clients := make([]*websocket.Conn, 6)
	for index := range clients {
		playerID := "player-" + string(rune('1'+index))
		options := &websocket.DialOptions{HTTPHeader: http.Header{"Authorization": []string{"Bearer test-" + playerID}}}
		connection, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws?matchId="+roomID+"&playerId="+playerID, options)
		if err != nil {
			t.Fatal(err)
		}
		clients[index] = connection
		_, _, _ = connection.Read(ctx)
	}
	defer func() {
		for _, client := range clients {
			client.CloseNow()
		}
	}()
	revision := uint64(0)
	send := func(index int, commandType, id string, sequence uint64, payload map[string]any) {
		playerID := "player-" + string(rune('1'+index))
		message := map[string]any{"version": "1.0.0", "messageId": id, "command": engine.Command{CommandID: id, MatchID: roomID, PlayerID: playerID, ExpectedRevision: revision, CommandType: commandType, Payload: payload, ClientTimestamp: time.Unix(1, 0).UTC(), ClientSequence: sequence}}
		data, _ := json.Marshal(message)
		if err := clients[index].Write(ctx, websocket.MessageText, data); err != nil {
			t.Fatal(err)
		}
		for {
			_, data, err = clients[index].Read(ctx)
			if err != nil {
				t.Fatal(err)
			}
			var envelope map[string]any
			_ = json.Unmarshal(data, &envelope)
			if envelope["type"] == "match.privateState" {
				revision++
				return
			}
		}
	}
	for index := 0; index < 6; index++ {
		send(index, engine.CommandJoinRoom, "join-"+string(rune('1'+index)), 1, nil)
		send(index, engine.CommandSetReady, "ready-"+string(rune('1'+index)), 2, map[string]any{"ready": true})
	}
	send(0, engine.CommandStartMatch, "start", 3, nil)
	for index := 0; index < 6; index++ {
		sequence := uint64(3)
		if index == 0 {
			sequence = 4
		}
		send(index, engine.CommandAcknowledgeRole, "ack-"+string(rune('1'+index)), sequence, nil)
	}
	room, _ := registry.Get(roomID)
	state, err := room.State(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if string(state.Phase) != "CHAPTER_START" {
		t.Fatalf("phase=%s", state.Phase)
	}
	private, ok := engine.ProjectPrivate(state, "player-1")
	if !ok || private.RoleID == "" || !private.RoleAcknowledged {
		t.Fatal("reconnect projection incomplete")
	}
}

func TestWebSocketRejectsUnauthenticatedAndMalformedMessages(t *testing.T) {
	registry := rooms.NewRegistry()
	defer registry.Close()
	server := httptest.NewServer(newHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), registry))
	defer server.Close()
	response, err := http.Post(server.URL+"/v1/rooms", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	var created map[string]string
	_ = json.NewDecoder(response.Body).Decode(&created)
	response.Body.Close()
	roomID := created["roomId"]
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, unauthorized, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws?matchId="+roomID+"&playerId=player-1", nil)
	if err == nil || unauthorized == nil || unauthorized.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized upgrade, got response=%v err=%v", unauthorized, err)
	}
	options := &websocket.DialOptions{HTTPHeader: http.Header{"Authorization": []string{"Bearer test-player-1"}}}
	connection, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws?matchId="+roomID+"&playerId=player-1", options)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseNow()
	_, _, _ = connection.Read(ctx)
	if err := connection.Write(ctx, websocket.MessageText, []byte(`{"version":"9.0.0","messageId":"bad"}`)); err != nil {
		t.Fatal(err)
	}
	_, data, err := connection.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("INVALID_ENVELOPE")) {
		t.Fatalf("unexpected protocol response: %s", data)
	}
}
