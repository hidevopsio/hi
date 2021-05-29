package service

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
)

// StatusHandler is the websocket handler
type StatusHandler struct {
	at.ContextAware
	connection *websocket.Connection
}

func newStatusHandler(connection *websocket.Connection) *StatusHandler {
	h := &StatusHandler{connection: connection}
	return h
}

func init() {
	app.Register(newStatusHandler)
}

// OnMessage is the websocket message handler
func (h *StatusHandler) OnMessage(data []byte) {
	message := string(data)
	log.Debugf("client: %v", message)

	h.connection.EmitMessage([]byte(fmt.Sprintf("Status: Up")))

}

// OnDisconnect is the websocket disconnection handler
func (h *StatusHandler) OnDisconnect() {
	log.Debugf("Connection with ID: %v has been disconnected!", h.connection.ID())
}



// OnPing is the websocket ping handler
func (h *StatusHandler) OnPing() {
	log.Debugf("Connection with ID: %v has been pinged!", h.connection.ID())
}


// OnPong is the websocket pong handler
func (h *StatusHandler) OnPong() {
	log.Debugf("Connection with ID: %v has been ponged!", h.connection.ID())
}