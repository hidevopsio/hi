package logging

import (
	"net/http"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
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
func (m *loggingMiddleware) Logging( _ struct{at.MiddlewareHandler `value:"/" `}, ctx context.Context) {
	ann := annotation.GetAnnotation(ctx.Annotations(), at.Operation{})
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
func (m *loggingMiddleware) PostLogging( _ struct{at.MiddlewarePostHandler `value:"/user/query" `}, ctx context.Context) {
	responses := ctx.GetResponses()
	var baseResponseInfo model.BaseResponseInfo
	var err error

	for _, resp := range responses {
		log.Debug(resp)
		if reflector.HasEmbeddedFieldType(resp, model.BaseResponseInfo{}) {
			respVal := reflector.GetFieldValue(resp, "BaseResponseInfo")
			if respVal.IsValid() {
				r := respVal.Interface()
				baseResponseInfo = r.(model.BaseResponseInfo)
			}
		}
		if resp != nil {
			switch resp.(type) {
			case error:
				log.Debug(resp)
				err = resp.(error)
				log.Warn(err)
			}
		}
	}
	if err == nil {
		log.Debugf("%v, %+v", ctx.Path(), baseResponseInfo)
	} else {
		baseResponseInfo.SetCode(http.StatusInternalServerError)
		baseResponseInfo.SetMessage("Internal server error")
		log.Debugf("%v: %v, %+v", ctx.Path(), err, baseResponseInfo)
	}

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

