package context

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"net/http"
)

// ExtendedContext extended context
type ExtendedContext interface {
	//RequestBody(data interface{}) error
	//RequestForm(data interface{}) error
	//RequestParams(request interface{}) error
	ResponseString(s string)
	ResponseBody(message string, data interface{})
	ResponseError(message string, code int)
	WrapHandler(h http.Handler)
	SetAnnotations(ann interface{})
	SetURLParam(name, value string)
	Annotations() interface{}
	AddResponse(response interface{})
	GetResponses() (responses map[string]interface{})
	GetResponse(object interface{}) (response interface{}, ok bool)
//StaticResource(system http.FileSystem)
}

// Context is the interface of web app context
type Context interface {
	context.Context
	ExtendedContext
}

// Handler is the handler func type (for Middleware)
type Handler func(Context)

// NewHandler will convert iris handler to our handler of func(*Context),
// in order to be compatible with the HTTP API.
func NewHandler(h iris.Handler) Handler {
	return func(ctx Context) {
		h(ctx.(context.Context))
	}
}
