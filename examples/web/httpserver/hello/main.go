package main

// import web starter from hiboot
import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
)

// main function
func main() {
	web.NewApplication().
		SetProperty(app.ProfilesInclude, logging.Profile, actuator.Profile).
		SetProperty(web.ViewEnabled, true).
		Run()
}
