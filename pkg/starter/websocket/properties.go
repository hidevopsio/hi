package websocket

import "github.com/hidevopsio/hiboot/pkg/at"

type properties struct {
	at.ConfigurationProperties `value:"websocket"`
	at.AutoWired

	// EvtMessagePrefix is the prefix of the underline websocket events that are being established under the hoods.
	// This prefix is visible only to the javascript side (code) and it has nothing to do
	// with the message that the end-user receives.
	// Do not change it unless it is absolutely necessary.
	//
	// If empty then defaults to []byte("websocket:").
	EvtMessagePrefix string `default:"websocket:"`
	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout int64	`default:"0"`
	// WriteTimeout time allowed to write a message to the connection.
	// 0 means no timeout.
	// Default value is 0
	WriteTimeout int64	`default:"0"`
	// ReadTimeout time allowed to read a message from the connection.
	// 0 means no timeout.
	// Default value is 0
	ReadTimeout int64	`default:"0"`
	// PongTimeout allowed to read the next pong message from the connection.
	// Default value is 60 * time.Second
	PongTimeout int64	`default:"60"`
	// PingPeriod send ping messages to the connection within this period. Must be less than PongTimeout.
	// Default value is 60 *time.Second
	PingPeriod int64	`default:"60"`
	// MaxMessageSize max message size allowed from connection.
	// Default value is 1024
	MaxMessageSize int64	`default:"1024"`
	// BinaryMessages set it to true in order to denotes binary data messages instead of utf-8 text
	// compatible if you wanna use the Connection's EmitMessage to send a custom binary data to the client, like a native server-client communication.
	// Default value is false
	BinaryMessages bool

	// Ping enable/disable ping, disabled by default
	Ping bool `default:"false"`

	ReadBufferSize  int    `default:"1024"`
	WriteBufferSize int    `default:"1024"`
	Javascript      string `default:"/websocket/websocket.js"`
}
