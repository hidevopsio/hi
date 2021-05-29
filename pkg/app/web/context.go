// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"net/http"
	"sync"

	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/kataras/iris"

	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/utils/mapstruct"
	"github.com/hidevopsio/hiboot/pkg/utils/validator"
	ctx "github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/i18n"
)

// Context Create your own custom Context, put any fields you wanna need.
type Context struct {
	iris.Context
	ann interface{}
	responses cmap.ConcurrentMap
}

//NewContext constructor of context.Context
func NewContext(app ctx.Application) context.Context {
	return &Context{
		Context: ctx.NewContext(app),
	}
}

var contextPool = sync.Pool{New: func() interface{} {
	return &Context{}
}}

func acquire(original iris.Context) *Context {
	c := contextPool.Get().(*Context)
	switch original.(type) {
	case *Context:
		newCtx := original.(*Context)
		c.Context = newCtx.Context
		c.responses = newCtx.responses
		c.ann = newCtx.ann
	default:
		c.Context = original // set the context to the original one in order to have access to iris's implementation.
	}
	return c
}

func release(c *Context) {
	contextPool.Put(c)
}

// Handler will convert our handler of func(*Context) to an iris Handler,
// in order to be compatible with the HTTP API.
func Handler(h func(context.Context)) iris.Handler {
	return func(original iris.Context) {
		c := acquire(original)
		h(c)
		release(c)
	}
}

//// WrapHandler is a helper function for wrapping http.Handler and returns a web middleware.
//func WrapHandler(h http.Handler) iris.Handler {
//	return Handler(func(c context.Context) {
//		h.ServeHTTP(c.ResponseWriter(), c.Request())
//	})
//}

// Next The second one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the chain of handlers via this "*Context".
func (c *Context) Next() {
	ctx.Next(c)
}

//// StaticResource is a helper function for wrapping http.Handler
//func (c *Context) StaticResource(system http.FileSystem) {
//	path := c.GetCurrentRoute().Path()
//	path = strings.Replace(path, "*", "", -1)
//	c.WrapHandler(http.StripPrefix(path, http.FileServer(system)))
//}

// WrapHandler is a helper function for wrapping http.Handler
func (c *Context) WrapHandler(h http.Handler)  {
	h.ServeHTTP(c.ResponseWriter(), c.Request())
}

// HTML Override any context's method you want...
// [...]
func (c *Context) HTML(htmlContents string) (int, error) {
	c.Application().Logger().Infof("Executing .HTML function from Context")

	c.ContentType("text/html")
	return c.WriteString(htmlContents)
}

// handle i18n
func (c *Context) translate(message string) string {

	message = i18n.Translate(c, message)

	return message
}

// Translate override base context method Translate to return format if i18n is not enabled
func (c *Context) Translate(format string, args ...interface{}) string {

	msg := c.Context.Translate(format, args...)

	if msg == "" {
		msg = format
	}

	return msg
}

// ResponseString set response
func (c *Context) ResponseString(data string) {
	//log.Infof("+++ %v : %v", c.Path(), c.translate(data))
	_, _ = c.WriteString(c.translate(data))
}

// ResponseBody set response
func (c *Context) ResponseBody(message string, data interface{}) {

	// TODO: check if data is a string, should we translate it?
	response := new(model.BaseResponse)
	response.SetCode(c.GetStatusCode())
	response.SetMessage(c.translate(message))
	response.SetData(data)

	c.JSON(response)
}

// ResponseError response with error
func (c *Context) ResponseError(message string, code int) {

	response := new(model.BaseResponse)
	response.SetCode(code)
	response.SetMessage(c.translate(message))
	if c.ResponseWriter() != nil {
		c.StatusCode(code)
		c.JSON(response)
	}
}

// SetAnnotations
func (c *Context) SetAnnotations(ann interface{})  {
	c.ann = ann
}

// Annotations
func (c *Context) Annotations() interface{} {
	return c.ann
}

// AddURLParam
func (c *Context) SetURLParam(name, value string) {
	q := c.Request().URL.Query()
	if q[name] != nil {
		q.Set(name, value)
	} else {
		q.Add(name, value)
	}
	c.Request().URL.RawQuery = q.Encode()
}


// AddResponse add response to a alice
func (c *Context) AddResponse(response interface{}) {
	if c.responses == nil {
		c.responses = cmap.New()
	}
	// TODO: do we need the index of the response value?
	name, object := factory.ParseParams(response)
	switch response.(type) {
	case error:
		c.responses.Set("error", object)
	default:
		if name == "" {
			// assume that name == "" means it is nil error
			c.responses.Set("error", response)
		} else {
			c.responses.Set(name, response)
		}
	}

	return
}

// GetResponses get all responses as a slice
func (c *Context) GetResponses() (responses map[string]interface{}) {
	if c.responses != nil {
		responses = c.responses.Items()
	}
	return
}

// GetResponse get specific response from a slice
func (c *Context) GetResponse(object interface{}) (response interface{}, ok bool) {
	name, _ := factory.ParseParams(object)
	if c.responses != nil {
		response, ok = c.responses.Get(name)
	}
	return
}

// RequestEx get RequestBody
func requestEx(c context.Context, data interface{}, cb func() error) error {
	if cb != nil {
		err := cb()
		if err != nil {
			c.ResponseError(err.Error(), http.StatusInternalServerError)
			return err
		}

		err = validator.Validate.Struct(data)
		if err != nil {
			c.ResponseError(err.Error(), http.StatusBadRequest)
			return err
		}
	}
	return nil
}

// RequestBody get RequestBody
func RequestBody(c context.Context, data interface{}) error {

	return requestEx(c, data, func() error {
		return c.ReadJSON(data)
	})
}

// RequestForm get RequestFrom
func RequestForm(c context.Context, data interface{}) error {

	return requestEx(c, data, func() error {
		return c.ReadForm(data)
	})
}

// RequestParams get RequestParams
func RequestParams(c context.Context, data interface{}) error {

	return requestEx(c, data, func() error {

		values := c.URLParams()
		if len(values) != 0 {
			return mapstruct.Decode(data, values)
		}
		return nil
	})
}
