// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package websocket provides web socket auto configuration for web/cli application
package websocket

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket/ws"
	"github.com/hidevopsio/hiboot/pkg/utils/copier"
	"time"
)

const (
	// Profile websocket profile name
	Profile = "websocket"
	// All is the string which the Emitter use to send a message to all.
	All = ""
	// Broadcast is the string which the Emitter use to send a message to all except this connection.
	Broadcast = ";to;all;except;me;"
)

type configuration struct {
	// embedded annotation at.AutoConfiguration
	at.AutoConfiguration

	Properties *properties
}

type Server struct {
	*websocket.Server
}

func newConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.Register(newConfiguration)
}

// Server websocket server
func (c *configuration) Server() *Server {
	var cfg websocket.Config
	_ = copier.Copy(&cfg, &c.Properties)
	p := c.Properties
	s := websocket.New(websocket.Config{
		Ping:             p.Ping,
		EvtMessagePrefix: []byte(p.EvtMessagePrefix),
		HandshakeTimeout: time.Duration(p.HandshakeTimeout) * time.Second,
		WriteTimeout:     time.Duration(p.WriteTimeout) * time.Second,
		ReadTimeout:      time.Duration(p.ReadTimeout) * time.Second,
		PongTimeout:      time.Duration(p.PongTimeout) * time.Second,
		PingPeriod:       time.Duration(p.PingPeriod) * time.Second,
		MaxMessageSize:   p.MaxMessageSize,
		BinaryMessages:   p.BinaryMessages,
		ReadBufferSize:   p.ReadBufferSize,
		WriteBufferSize:  p.WriteBufferSize,
	})

	return &Server{
		Server: s,
	}
}

// Connection websocket connection for runtime dependency injection
func (c *configuration) Connection(ctx context.Context, server *Server) *Connection {
	conn := newConnection(server.Upgrade(ctx))
	conn.SetValue("context", ctx)
	return conn
}

// Register is function that register handler
func (c *configuration) Register() Register {
	return registerHandler
}
