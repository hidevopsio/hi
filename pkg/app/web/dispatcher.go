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
	"fmt"
	"github.com/fatih/camelcase"
	"github.com/kataras/iris/v12"
	"github.com/rakyll/statik/fs"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/app/web/webutils"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/system"
	"hidevops.io/hiboot/pkg/utils/reflector"
	"hidevops.io/hiboot/pkg/utils/str"
	"net/http"
	"path"
	"reflect"
	"strings"
)

var httpMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

const (
	Any             = "ANY"
	RequestMapping  = "RequestMapping"
	ContextPathRoot = "/"
	UrlSep          = "/"
)

type Dispatcher struct {
	webApp *webApp
	// inject context aware dependencies
	configurableFactory factory.ConfigurableFactory
	SystemApp           *system.App
	SystemServer        *system.Server

	ContextPathFormat string `value:"${server.context_path_format}" `

	methodSubscribers []*factory.MetaData
}

type requestMapping struct {
	Method string
	Value  string
}

type injectableMethod struct {
	method              *reflect.Method
	annotations         *annotation.Annotations
	handler             iris.Handler
	hasMethodAnnotation bool
	requestMapping      *requestMapping
}

type injectableObject struct {
	object      interface{}
	name        string
	pkgPath     string
	pathPrefix  string
	before      *injectableMethod
	after       *injectableMethod
	methods     []*injectableMethod
	annotations *annotation.Annotations
}

type Annotations struct {
	Fields []*annotation.Field
	Object interface{}
	Value  reflect.Value
}

func newDispatcher(webApp *webApp, configurableFactory factory.ConfigurableFactory) *Dispatcher {
	d := &Dispatcher{
		webApp:              webApp,
		configurableFactory: configurableFactory,
	}
	return d
}

func init() {
	app.Register(newDispatcher)
}

func (d *Dispatcher) parseAnnotation(object interface{}, method *reflect.Method) (ma *annotation.Annotations) {
	numIn := method.Type.NumIn()
	inputs := make([]reflect.Value, numIn)
	inputs[0] = reflect.ValueOf(object)
	for n := 1; n < numIn; n++ {
		typ := method.Type.In(n)
		iTyp := reflector.IndirectType(typ)
		if iTyp.Name() == "" && typ.Kind() == reflect.Struct {
			av := reflect.New(typ)
			avo := av.Interface()
			ma = annotation.GetAnnotations(avo)
			if ma != nil {
				_ = d.configurableFactory.InjectIntoObject(avo)
				break
			}
		}
	}
	return
}

// TODO: controller method handler, middleware handler, and event handler, they all use the same way to achieve handler dispatch.
func (d *Dispatcher) parseMiddleware(m *factory.MetaData) (middleware *injectableObject) {
	middleware = new(injectableObject)

	mwi := reflect.ValueOf(m.Instance)
	middleware.object = mwi.Interface()

	// set annotations
	annotations := annotation.GetAnnotations(middleware.object)
	middleware.annotations = annotations

	mwType := mwi.Type()
	numOfMethod := mwi.NumMethod()
	for mi := 0; mi < numOfMethod; mi++ {
		methodHandler := new(injectableMethod)
		method := mwType.Method(mi)
		methodHandler.method = &method
		ma := d.parseAnnotation(mi, &method)
		methodHandler.annotations = ma
		mwh := annotation.FilterIn(ma, at.MiddlewareHandler{})
		if len(mwh) > 0 {
			methodHandler.hasMethodAnnotation = true
			hdl := newHandler(d.configurableFactory, middleware, methodHandler, at.MiddlewareHandler{})
			methodHandler.handler = Handler(func(c context.Context) {
				hdl.call(c)
			})
			middleware.methods = append(middleware.methods, methodHandler)
		}
	}
	return
}

func (d *Dispatcher) parseRestController(ctl *factory.MetaData) (restController *injectableObject) {

	restController = new(injectableObject)

	c := ctl.Instance
	field := reflect.ValueOf(c)

	fieldType := field.Type()
	//log.Debug("fieldType: ", fieldType)
	ift := fieldType.Elem()
	fieldName := ift.Name()
	restController.pkgPath = ift.PkgPath()
	//log.Debug("fieldName: ", fieldName)

	controller := field.Interface()
	restController.object = controller
	//log.Debug("controller: ", controller)
	annotations := annotation.GetAnnotations(controller)
	restController.annotations = annotations

	// get context mapping
	var customizedControllerPath bool
	pathPrefix := d.SystemServer.ContextPath
	af := annotation.GetAnnotation(controller, at.RequestMapping{})
	if af != nil {
		customizedControllerPath = true
		ann := af.Field.Value.Interface().(at.RequestMapping)
		pathPrefix = path.Join(pathPrefix, ann.AtValue)
	}

	// parse method
	fieldNames := camelcase.Split(fieldName)
	controllerName := ""
	if len(fieldNames) >= 2 {
		controllerName = strings.Replace(fieldName, fieldNames[len(fieldNames)-1], "", 1)
		controllerName = str.LowerFirst(controllerName)
	}
	//log.Debug("controllerName: ", controllerName)
	// use controller's prefix as context mapping
	if !customizedControllerPath {
		cn := controllerName
		switch d.ContextPathFormat {
		case app.ContextPathFormatKebab:
			cn = str.ToKebab(controllerName)
		case app.ContextPathFormatSnake:
			cn = str.ToSnake(controllerName)
		case app.ContextPathFormatCamel:
			cn = str.ToCamel(controllerName)
		case app.ContextPathFormatLowerCamel:
			cn = str.ToLowerCamel(controllerName)
		}
		contextPath := d.SystemServer.ContextPath
		//if contextPath == ContextPathRoot {
		//	contextPath = ""
		//}
		pathPrefix = fmt.Sprintf("%v/%v", contextPath, cn)
	}
	restController.pathPrefix = pathPrefix
	restController.name = fieldName

	numOfMethod := field.NumMethod()
	//log.Debug("methods: ", numOfMethod)

	// find before, after method
	before, ok := fieldType.MethodByName(beforeMethod)
	if ok {
		restMethod := new(injectableMethod)
		restMethod.method = &before
		restController.before = restMethod
	}

	after, ok := fieldType.MethodByName(afterMethod)
	if ok {
		restMethod := new(injectableMethod)
		restMethod.method = &after
		restController.after = restMethod
	}

	var methods []*injectableMethod
	for mi := 0; mi < numOfMethod; mi++ {
		restMethod := new(injectableMethod)

		method := fieldType.Method(mi)
		methodName := method.Name
		restMethod.annotations = d.parseAnnotation(controller, &method)
		restMethod.method = &method

		reqMap := new(requestMapping)
		m, p := webutils.GetHttpMethod(restMethod.annotations)
		if m != "" {
			restMethod.hasMethodAnnotation = true
			reqMap.Method, reqMap.Value = m, p
		}

		beforeMethod := annotation.FilterIn(restMethod.annotations, at.BeforeMethod{})
		if len(beforeMethod) > 0 {
			restMethod.method = &method
			restController.before = restMethod
			continue
		}
		afterMethod := annotation.FilterIn(restMethod.annotations, at.AfterMethod{})
		if len(afterMethod) > 0 {
			restMethod.method = &method
			restController.after = restMethod
			continue
		}

		if !restMethod.hasMethodAnnotation {
			ctxMap := camelcase.Split(methodName)
			reqMap.Method = strings.ToUpper(ctxMap[0])
			var apiPath string
			if len(ctxMap) > 2 && ctxMap[1] == "By" {
				for _, pathParam := range ctxMap[2:] {
					lpp := strings.ToLower(pathParam)
					apiPath = apiPath + pathSep + lpp + pathSep + "{" + lpp + "}"
				}
			} else {
				apiPath = strings.Replace(methodName, ctxMap[0], "", 1)
				apiPath = pathSep + str.LowerFirst(apiPath)
			}
			reqMap.Value = apiPath
		}
		restMethod.requestMapping = reqMap

		hasAnyMethod := reqMap.Method == Any
		hasRegularMethod := str.InSlice(reqMap.Method, httpMethods)
		foundMethod := hasAnyMethod || hasRegularMethod
		if foundMethod {
			methods = append(methods, restMethod)
		}

	}
	restController.methods = methods
	return
}

func (d *Dispatcher) useMiddleware(mw []*annotation.Annotation, mwMth []*annotation.Annotation, ctl []*annotation.Annotation, ctlMth []*annotation.Annotation) (matched bool) {
	var atMwTyp reflect.Type

	// (a || b) && (c || d)
	// 0 0 x x => 1
	if len(mw) > 0 {
		atMwTyp = mw[0].Field.StructField.Type
	}
	if atMwTyp == nil && len(mwMth) > 0 {
		atMwTyp = mwMth[0].Field.StructField.Type
	}
	if atMwTyp == nil {
		matched = true
		return
	}

	if len(ctl) > 0 {
		for _, c := range ctl {
			if atMwTyp == c.Field.StructField.Type {
				matched = true
				return
			}
		}
	}

	if len(ctlMth) > 0 {
		for _, cm := range ctlMth {
			if atMwTyp == cm.Field.StructField.Type {
				matched = true
				return
			}
		}
	}

	// did not match
	// 0 1 0 0 => 0
	// 1 0 0 0 => 0
	// 1 1 0 0 => 0

	return
}

//TODO: scan apis and params to generate swagger api automatically by include swagger starter
func (d *Dispatcher) register(controllers []*factory.MetaData, middleware []*factory.MetaData) (err error) {

	d.methodSubscribers = d.configurableFactory.GetInstances(at.HttpMethodSubscriber{})

	var mws []*injectableObject
	for _, m := range middleware {
		mw := d.parseMiddleware(m)
		if mw != nil {
			mws = append(mws, mw)
		}
	}

	log.Debug("register rest controller")
	for _, ctl := range controllers {
		// get and parse all controller methods
		restController := d.parseRestController(ctl)

		var party iris.Party
		if restController.before != nil {
			hdl := newHandler(d.configurableFactory, restController, restController.before, at.BeforeMethod{})
			party = d.webApp.Party(restController.pathPrefix, Handler(func(c context.Context) {
				hdl.call(c)
			}))
		} else {
			party = d.webApp.Party(restController.pathPrefix)
		}

		if restController.after != nil {
			hdl := newHandler(d.configurableFactory, restController, restController.after, at.AfterMethod{})
			party.Done(Handler(func(c context.Context) {
				hdl.call(c)
			}))
		}

		// bind regular http method handlers with router
		atCtl := annotation.FilterIn(restController.annotations, at.UseMiddleware{})
		for _, m := range restController.methods {
			var handlers []iris.Handler

			atCtlMth := annotation.FilterIn(m.annotations, at.UseMiddleware{})

			// 1. pass all annotations to registered starter for further implementations, e.g. swagger

			// 2. get all handlers from middleware, then use middleware
			// first check if controller annotated to use at.Middleware
			// second check if method annotated to use at.MiddlewareHandler
			// handlers = append(handlers, middleware...)

			// set matched to true by default
			if len(mws) > 0 {
				for _, mw := range mws {
					atMw := annotation.FilterIn(mw.annotations, at.UseMiddleware{})
					for _, mth := range mw.methods {
						atMwMth := annotation.FilterIn(mth.annotations, at.UseMiddleware{})
						// check if middleware is used
						// else skip to append this middleware
						useMiddleware := d.useMiddleware(atMw, atMwMth, atCtl, atCtlMth)
						if useMiddleware {
							handlers = append(handlers, mth.handler)
						}
					}
				}
			}

			// 3. finally, handle all method handlers
			d.handleControllerMethod(restController, m, party, handlers)
		}
	}
	return
}

func (d *Dispatcher) handleControllerMethod(restController *injectableObject, m *injectableMethod, party iris.Party, handlers []iris.Handler) {
	// 3. create new handler for rest controller method
	hdl := newHandler(d.configurableFactory, restController, m, at.HttpMethod{})

	var h iris.Handler
	atFileServer := annotation.GetAnnotation(m.annotations, at.FileServer{})
	if atFileServer != nil {
		afs := atFileServer.Field.Value.Interface().(at.FileServer)
		path := restController.pathPrefix + afs.AtValue
		h = Handler(func(c context.Context) {
			// call controller method first
			hdl.call(c)

			// serve static resource
			f, err := fs.New()
			if err == nil {
				c.WrapHandler(http.StripPrefix(path, http.FileServer(f)))
			}

			// next
			c.Next()
		})
	} else {
		h = Handler(func(c context.Context) {
			hdl.call(c)
			c.Next()
		})
	}
	handlers = append(handlers, h)

	if m.requestMapping.Method == Any {
		party.Any(m.requestMapping.Value, handlers...)
	} else {
		route := party.Handle(m.requestMapping.Method, m.requestMapping.Value, handlers...)
		route.MainHandlerName = fmt.Sprintf("%s/%s.%s", restController.pkgPath, restController.name, m.method.Name)
	}

	// publish to subscriber
	for _, ms := range d.methodSubscribers {
		if reflector.Implements(ms.Instance, new(HttpMethodSubscriber)) {
			subscriber := ms.Instance.(HttpMethodSubscriber)
			subscriber.Subscribe(restController.annotations, m.annotations)
		}
	}
}
