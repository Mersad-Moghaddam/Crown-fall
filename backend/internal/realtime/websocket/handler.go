package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"crownfall/backend/internal/game/engine"
	"crownfall/backend/internal/realtime/rooms"
	ws "github.com/coder/websocket"
)

const (
	ProtocolVersion       = "1.0.0"
	maximumMessageBytes   = 16 << 10
	outboundQueueCapacity = 32
)

type ClientEnvelope struct {
	Version   string         `json:"version"`
	MessageID string         `json:"messageId"`
	Command   engine.Command `json:"command"`
}
type ServerEnvelope struct {
	Version    string    `json:"version"`
	MessageID  string    `json:"messageId"`
	MatchID    string    `json:"matchId,omitempty"`
	Sequence   uint64    `json:"sequence,omitempty"`
	ServerTime time.Time `json:"serverTime"`
	Type       string    `json:"type"`
	Payload    any       `json:"payload"`
}
type ProtocolError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type connection struct {
	socket   *ws.Conn
	outbound chan []byte
	cancel   context.CancelFunc
}
type Server struct {
	registry    *rooms.Registry
	mu          sync.Mutex
	connections map[string]map[string]*connection
}

func NewServer(registry *rooms.Registry) *Server {
	return &Server{registry: registry, connections: make(map[string]map[string]*connection)}
}

func (server *Server) Handler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		matchID, playerID := request.URL.Query().Get("matchId"), request.URL.Query().Get("playerId")
		if matchID == "" || playerID == "" || request.Header.Get("Authorization") != "Bearer test-"+playerID {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}
		room, exists := server.registry.Get(matchID)
		if !exists {
			http.Error(writer, "room not found", http.StatusNotFound)
			return
		}
		socket, err := ws.Accept(writer, request, &ws.AcceptOptions{OriginPatterns: []string{"localhost:*", "127.0.0.1:*"}})
		if err != nil {
			return
		}
		socket.SetReadLimit(maximumMessageBytes)
		ctx, cancel := context.WithCancel(request.Context())
		client := &connection{socket: socket, outbound: make(chan []byte, outboundQueueCapacity), cancel: cancel}
		server.register(matchID, playerID, client)
		defer func() { server.unregister(matchID, playerID, client); cancel(); socket.CloseNow() }()
		go server.writeLoop(ctx, client)

		state, err := room.State(ctx)
		if err == nil {
			if view, member := engine.ProjectResync(state, playerID); member {
				server.send(client, envelope("connection.resynced", "connection", matchID, state.EventSequence, view))
			} else {
				server.send(client, envelope("connection.accepted", "connection", matchID, state.EventSequence, map[string]any{"phase": state.Phase, "revision": state.Revision}))
			}
		}
		for {
			_, data, err := socket.Read(ctx)
			if err != nil {
				return
			}
			var message ClientEnvelope
			decoderError := json.Unmarshal(data, &message)
			if decoderError != nil || message.Version != ProtocolVersion || message.MessageID == "" {
				server.sendError(client, matchID, message.MessageID, "INVALID_ENVELOPE", "malformed envelope or unsupported protocol version")
				continue
			}
			if message.Command.PlayerID != playerID || message.Command.MatchID != matchID {
				server.sendError(client, matchID, message.MessageID, "IDENTITY_MISMATCH", "authenticated identity or match does not match command")
				continue
			}
			result, err := room.Handle(ctx, message.Command)
			if err != nil {
				server.sendError(client, matchID, message.MessageID, errorCode(err), err.Error())
				continue
			}
			public := engine.ProjectPublic(result.State)
			server.broadcastPublic(matchID, envelope("match.publicState", message.MessageID, matchID, result.State.EventSequence, public))
			if private, ok := engine.ProjectPrivate(result.State, playerID); ok {
				server.send(client, envelope("match.privateState", message.MessageID, matchID, result.State.EventSequence, private))
			}
			for recipient, events := range result.PrivateEvents {
				if target := server.get(matchID, recipient); target != nil {
					server.send(target, envelope("match.privateEvents", message.MessageID, matchID, result.State.EventSequence, events))
				}
			}
		}
	}
}

func envelope(eventType, messageID, matchID string, sequence uint64, payload any) ServerEnvelope {
	return ServerEnvelope{Version: ProtocolVersion, MessageID: messageID, MatchID: matchID, Sequence: sequence, ServerTime: time.Now().UTC(), Type: eventType, Payload: payload}
}
func (server *Server) sendError(client *connection, matchID, messageID, code, message string) {
	server.send(client, envelope("protocol.error", messageID, matchID, 0, ProtocolError{Code: code, Message: message}))
}
func (server *Server) send(client *connection, message ServerEnvelope) {
	data, _ := json.Marshal(message)
	select {
	case client.outbound <- data:
	default:
		client.cancel()
	}
}
func (server *Server) writeLoop(ctx context.Context, client *connection) {
	for {
		select {
		case data := <-client.outbound:
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := client.socket.Write(writeCtx, ws.MessageText, data)
			cancel()
			if err != nil {
				client.cancel()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
func (server *Server) register(matchID, playerID string, client *connection) {
	server.mu.Lock()
	defer server.mu.Unlock()
	if server.connections[matchID] == nil {
		server.connections[matchID] = make(map[string]*connection)
	}
	if old := server.connections[matchID][playerID]; old != nil {
		old.cancel()
	}
	server.connections[matchID][playerID] = client
}
func (server *Server) unregister(matchID, playerID string, client *connection) {
	server.mu.Lock()
	defer server.mu.Unlock()
	if server.connections[matchID][playerID] == client {
		delete(server.connections[matchID], playerID)
	}
	if len(server.connections[matchID]) == 0 {
		delete(server.connections, matchID)
	}
}
func (server *Server) get(matchID, playerID string) *connection {
	server.mu.Lock()
	defer server.mu.Unlock()
	return server.connections[matchID][playerID]
}
func (server *Server) broadcastPublic(matchID string, message ServerEnvelope) {
	server.mu.Lock()
	targets := make([]*connection, 0, len(server.connections[matchID]))
	for _, target := range server.connections[matchID] {
		targets = append(targets, target)
	}
	server.mu.Unlock()
	for _, target := range targets {
		server.send(target, message)
	}
}
func errorCode(err error) string {
	switch {
	case errors.Is(err, engine.ErrStaleRevision):
		return "STALE_REVISION"
	case errors.Is(err, engine.ErrHostOnly):
		return "HOST_ONLY"
	case errors.Is(err, engine.ErrNotMember):
		return "NOT_MEMBER"
	case errors.Is(err, rooms.ErrMailboxFull):
		return "MAILBOX_FULL"
	default:
		return "COMMAND_REJECTED"
	}
}
