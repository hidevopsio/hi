package main

import (
	_ "github.com/hidevopsio/hiboot/examples/web/websocket/controller"
	_ "github.com/hidevopsio/hiboot/examples/web/websocket/service"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/websocket"
)

func main() {
	web.NewApplication().
		SetProperty(app.ProfilesInclude, websocket.Profile).
		Run()
}
