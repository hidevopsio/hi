package main

import (
	_ "github.com/hidevopsio/hiboot/examples/web/middleware/controller"
	_ "github.com/hidevopsio/hiboot/examples/web/middleware/logging"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
)

// HiBoot main function
func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.ProfilesInclude, web.Profile, actuator.Profile, logging.Profile).
		Run()
}