package WsHandler

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/samber/lo"
)

type WsHub struct {
	Users map[int64]*Hub
}

func (w *WsHub) UnregisterWs(conn *websocket.Conn, userId int64) {
	w.Users[userId].Connections = lo.Filter(w.Users[userId].Connections, func(item *websocket.Conn, index int) bool {
		if item != conn {
			return true
		} else {
			return false
		}
	})
}

func (w *WsHub) RegisterWs(conn *websocket.Conn, userId int64) {

	if w == nil {
		return
	}

	_, exists := w.Users[userId]
	if exists == false {
		w.Users[userId] = &Hub{
			Connections: make([]*websocket.Conn, 0),
			Broadcast:   make(chan []byte, 1e5),
			Register:    make(chan []byte, 1e5),
			Unregister:  make(chan []byte, 1e5),
		}
		w.Users[userId].RegisterUserConnection(conn)
	} else {
		w.Users[userId].RegisterUserConnection(conn)
	}
}

func (w *WsHub) testBroadcastAmongAll(msg []byte) {
	for i := range w.Users {
		w.Users[i].BroadcastMessage(msg)
	}
}
