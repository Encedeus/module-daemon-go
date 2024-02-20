package module

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func StartEchoServer(m *Module) error {
	m.Echo.GET("/ws", hello)

	err := m.Echo.Start(fmt.Sprintf(":%v", m.Port))
	if err != nil {
		return err
	}

	return nil
}

func hello(c echo.Context) error {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, my dear friend!"))
		if err != nil {
			c.Logger().Error(err)
		}

		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("Websocket msg: %v", msg)
	}
}
