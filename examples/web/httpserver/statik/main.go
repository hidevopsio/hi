//go:generate statik -src=./static

package main

import (
	_ "github.com/hidevopsio/hiboot/examples/web/httpserver/statik/statik"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
)

type controller struct {
	at.RestController
	at.RequestMapping `value:"/public"`
}

func init() {
	app.Register(newStaticController)
}

func newStaticController() *controller {
	return &controller{}
}

// UI serve static resource via context StaticResource method
func (c *controller) UI(at struct{ at.GetMapping `value:"/ui/*"`; at.FileServer `value:"/ui"` }, ctx context.Context) {
	return
}

// UI serve static resource via context StaticResource method
func (c *controller) UIIndex(at struct{ at.GetMapping `value:"/ui"`; at.FileServer `value:"/ui"` }, ctx context.Context) {
	return
}

// Before run go build, run go generate.
// Then, run the main program and visit below urls:
// http://localhost:8080/public/ui
// http://localhost:8080/public/ui/hello.txt
// http://localhost:8080/public/ui/img/hiboot.png

func main() {
	web.NewApplication(newStaticController).
		SetProperty(app.ProfilesInclude, actuator.Profile, logging.Profile).
		Run()
}
