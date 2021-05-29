//go:generate statik -src=./dist

package main

import (
	_ "github.com/hidevopsio/hiboot/examples/web/swagger/greeting-server/controller"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/server"
	"github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"github.com/hidevopsio/hiboot/pkg/starter/swagger"
)

const (
	version = "1.0.3"
	basePath = "/api/v1/greeting-server"
)

func init() {
	app.Register(swagger.ApiInfoBuilder().
		ContactName("John Deng").
		ContactEmail("john.deng@outlook.com").
		ContactURL("https://hidevops.io").
		Title("HiBoot Swagger Demo Application - Greeting Server").
		Description("Greeting Server is an application that demonstrate the usage of Swagger Annotations"),
		// alternatively, you can set below properties by using SetProperty() in main, config/application.yml or program arguments to take advantage of HiBoot DI
		//Version(version).
		//Schemes("http", "https").
		//BasePath(basePath),
	)
}

//run http://localhost:8080/api/v1/greeting-server/swagger-ui to open swagger ui
func main() {
	web.NewApplication().
		SetProperty(app.ProfilesInclude, actuator.Profile, logging.Profile, swagger.Profile).
		SetProperty(app.Version, version).
		SetProperty(server.Host, "localhost").
		SetProperty(server.Port, "8080").
		SetProperty(server.ContextPath, basePath).
		Run()
}
