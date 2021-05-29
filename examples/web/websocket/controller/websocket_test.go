package controller

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
	"net/http"
	"testing"
)

func TestWebSocketController(t *testing.T) {
	mockController := newWebsocketController(func(handler websocket.Handler, conn *websocket.Connection) {
		// For controller's unit testing, do nothing
		ctx := conn.GetValue("context").(context.Context)
		ctx.StatusCode(http.StatusOK)
	})

	testApp := web.NewTestApp(mockController).SetProperty(app.ProfilesInclude, websocket.Profile, logging.Profile).Run(t)
	assert.NotEqual(t, nil, testApp)

	testApp.Get("/websocket").Expect().Status(http.StatusServiceUnavailable)
	testApp.Get("/websocket/status").Expect().Status(http.StatusServiceUnavailable)
}
