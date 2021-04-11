package logging

import (
	"github.com/hidevopsio/hiboot/examples/web/middleware/controller"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type loggingMiddleware struct {
	at.Middleware
}

func newLoggingMiddleware() *loggingMiddleware {
	return &loggingMiddleware{}
}

func init() {
	app.Register(newLoggingMiddleware)
}

// Logging is the middleware handler,it support dependency injection, method annotation
// middleware handler can be annotated to specific purpose or general purpose
func (m *loggingMiddleware) Before( _ struct{at.BeforeHandler `value:"/" `}, ctx context.Context) {

	// TODO: test code only
	ctx.Values().Set("post-process", "true")

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

// PostLogging is the middleware post handler
func (m *loggingMiddleware) After( _ struct{at.AfterHandler `value:"/" `}, ctx context.Context) {
	res, ok := ctx.GetResponse(controller.UserResponse{})
	if ok {
		response := res.(*controller.UserResponse)
		log.Debugf("%+v", response)
		response.Data.Username = "changed by middleware"
		ctx.StatusCode(503)
		_, _ = ctx.JSON(response.Data)
	}
	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

