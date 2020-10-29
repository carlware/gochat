package rest

import (
	"github.com/carlware/gochat/dispatchers/websocket"
	"github.com/labstack/echo/v4"
)

// WS starts websocket for the client
func (ctr *controller) WS(c echo.Context) error {
	websocket.ServeWs(ctr.hub, c.Response(), c.Request())
	return nil
}
