package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

type ConnectionEvent struct {
	ProtocolVersion string    `json:"protocol_version"`
	Type            string    `json:"type"`
	ServerTime      time.Time `json:"server_time"`
}

func Handler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		connection, err := websocket.Accept(writer, request, &websocket.AcceptOptions{OriginPatterns: []string{"localhost:*"}})
		if err != nil {
			return
		}
		defer connection.CloseNow()
		payload, _ := json.Marshal(ConnectionEvent{ProtocolVersion: "1.0.0", Type: "connection.accepted", ServerTime: time.Now().UTC()})
		ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
		defer cancel()
		if err := connection.Write(ctx, websocket.MessageText, payload); err != nil {
			return
		}
		connection.Close(websocket.StatusNormalClosure, "initial scaffold complete")
	}
}
