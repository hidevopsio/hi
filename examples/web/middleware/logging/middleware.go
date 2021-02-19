package logging

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
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
func (m *loggingMiddleware) Logging( a struct{at.MiddlewareHandler `value:"/" `}, ctx context.Context) {
	ann := annotation.GetAnnotation(ctx.Annotations(), at.RequiresPermissions{})
	if ann != nil {
		va := ann.Field.Value.Interface().(at.RequiresPermissions)
		switch va.AtType {
		case "path":
			id := ctx.Params().Get(va.AtIn[0])
			log.Debugf("path id: %v", id)
		case "query":
			id := ctx.URLParams()[va.AtIn[0]]
			log.Debugf("query id: %v", id)
		case "query:pagination":
			log.Debugf("page number: %v, page size: %v", ctx.URLParam(va.AtIn[0]), ctx.URLParam(va.AtIn[1]))
			ctx.SetURLParam(va.AtOut[0], "in(2,4,6,8)")
			ctx.SetURLParam(va.AtOut[1], "50")
			ctx.Header(va.AtOut[0], "in(2,4,6,8)")
			ctx.Request().Header.Set(va.AtOut[0], "in(2,4,6,8)")
		}

		log.Infof("[auth middleware] %v - %v", ctx.GetCurrentRoute(), va.AtValues)
	}

	ann = annotation.GetAnnotation(ctx.Annotations(), at.Operation{})
	if ann != nil {
		va := ann.Field.Value.Interface().(at.Operation)
		log.Infof("[logging middleware] %v - %v", ctx.GetCurrentRoute(), va.AtDescription)
	} else {
		log.Infof("[logging middleware] %v", ctx.GetCurrentRoute())
	}

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

// PostLogging is the middleware post handler
func (m *loggingMiddleware) PostLogging( a struct{at.MiddlewarePostHandler `value:"/" `}, ctx context.Context) {
	ann := annotation.GetAnnotation(ctx.Annotations(), at.Operation{})
	if ann != nil {
		va := ann.Field.Value.Interface().(at.Operation)
		log.Infof("[post logging middleware] %v - %v - %v", ctx.GetCurrentRoute(), ctx.GetCurrentRoute(), va.AtDescription)
	} else {
		log.Infof("[post logging middleware] %v", ctx.GetCurrentRoute())
	}

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

