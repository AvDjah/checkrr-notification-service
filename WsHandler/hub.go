package WsHandler

import (
	"checkrr-notification-service/Helpers"
	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	Connections []*websocket.Conn
	Broadcast   chan []byte
	Register    chan []byte
	Unregister  chan []byte
}

func (h *Hub) RegisterUserConnection(conn *websocket.Conn) {
	h.Connections = append(h.Connections, conn)
}

func (h *Hub) BroadcastMessage(msg []byte) {
	for _, conn := range h.Connections {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			Helpers.Log(err, "Error Broadcasting message")

		}
	}
}
