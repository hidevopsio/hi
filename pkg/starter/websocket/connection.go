package websocket

import (
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket/ws"
)

// Connection is the websocket connection
type Connection struct {
	at.ContextAware
	websocket.Connection
}

func newConnection(conn websocket.Connection) *Connection {
	return &Connection{Connection: conn}
}
