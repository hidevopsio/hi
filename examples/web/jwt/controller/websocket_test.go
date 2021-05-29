package controller

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
	"net/http"
	"testing"
)

func TestWebSocketController(t *testing.T) {
	mockController := newWebsocketController(func(handler websocket.Handler, conn *websocket.Connection) {
		// For controller's unit testing, do nothing
		log.Debug("mock controller")
		ctx := conn.GetValue("context").(context.Context)
		ctx.StatusCode(http.StatusOK)
	})

	testApp := web.NewTestApp(mockController).Run(t)
	assert.NotEqual(t, nil, testApp)

	token := genJwtToken(1000)

	testApp.Get("/websocket").
		WithHeader("Authorization", token).
		Expect().Status(http.StatusServiceUnavailable)

	testApp.Get("/websocket/status").
		WithHeader("Authorization", token).
		Expect().Status(http.StatusServiceUnavailable)
}
