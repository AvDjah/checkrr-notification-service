package Router

import (
	"checkrr-notification-service/Helpers"
	"checkrr-notification-service/WsHandler"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func New(ch chan []byte) *fiber.App {
	app := fiber.New()

	WsHub := WsHandler.WsHub{
		Users: make(map[int64]*WsHandler.Hub),
	}
	go RunBroadcaster(ch, &WsHub)
	AddRoutes(app, ch, &WsHub)

	return app
}

func RunBroadcaster(ch chan []byte, wsHub *WsHandler.WsHub) {
	for {
		select {
		case msg := <-ch:
			for val := range wsHub.Users {
				wsHub.Users[val].BroadcastMessage(msg)
			}
		}
	}
}

func ListenForClose(c *websocket.Conn, done chan string) {

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err) {
				fmt.Println("Closing Connection")
				done <- "Done"
			} else if websocket.IsUnexpectedCloseError(err) {
				fmt.Println("Unexpected Closing Connection")
				done <- "Done"
			}

		}
	}
}

func AddRoutes(app *fiber.App, ch chan []byte, wsHub *WsHandler.WsHub) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://gofiber.io, http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {

		userId := Helpers.GetUserIdFromJWTClaim(c)
		if userId == 0 {
			err := c.Close()
			if err != nil {
				Helpers.Log(err, "Closing WebSocket Connection")
				return
			}
		}

		forever := make(chan string, 2)

		go ListenForClose(c, forever)

		c.SetCloseHandler(func(code int, text string) error {
			fmt.Println("Closing Channel: ", code, text)
			forever <- "Close"
			return nil
		})

		wsHub.RegisterWs(c, userId)
		<-forever
		wsHub.UnregisterWs(c, userId)
		fmt.Println("Closed Connection Finna")
		//var (
		//	//mt  int
		//	//msg []byte
		//	err error
		//)
		//
		//done := make(chan string)
		//go func() {
		//	defer close(done)
		//	for {
		//		_, message, err := c.ReadMessage()
		//		if err != nil {
		//			log.Println("read:", err)
		//			break
		//		}
		//		fmt.Println("Received:", string(message))
		//		// Process the received message here
		//	}
		//}()
		//
		//for {
		//
		//	select {
		//	case msg := <-ch:
		//		fmt.Println("Written : ", string(msg))
		//		if err = c.WriteMessage(websocket.TextMessage, msg); err != nil {
		//			log.Println("write:", err)
		//			break
		//		}
		//	}
		//}

		//<-done
	}))
}
