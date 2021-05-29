package controller

import (
	"github.com/hidevopsio/hiboot/examples/web/jwt/service"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
)

type websocketController struct {
	at.JwtRestController

	register websocket.Register
}

func newWebsocketController(register websocket.Register) *websocketController {
	return &websocketController{register: register}
}

func init() {
	app.Register(newWebsocketController)
}

// Get GET /websocket
func (c *websocketController) Get(handler *service.CountHandler, connection *websocket.Connection) {
	c.register(handler, connection)
}

// GetStatus GET /websocket/status
func (c *websocketController) GetStatus(handler *service.StatusHandler, connection *websocket.Connection) {
	c.register(handler, connection)
}
