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
//go:generate statik -src=./static

// package web_test provides web uint tests
package web_test

import (
	"errors"
	"fmt"
	"sync"

	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	_ "github.com/hidevopsio/hiboot/pkg/app/web/statik"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/starter/locale"

	//_ "github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
	_ "github.com/hidevopsio/hiboot/pkg/starter/locale"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	_ "github.com/hidevopsio/hiboot/pkg/starter/logging"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"net/http"
	"testing"
	"time"
)

var mu sync.Mutex

type UserRequest struct {
	model.RequestBody
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type FooRequest struct {
	at.RequestBody
	Name string
}

type BarRequest struct {
	at.RequestParams
	Name string
}

type FoobarRequestForm struct {
	at.RequestForm
	Name string
}

type FoobarRequestParams struct {
	at.RequestParams
	Name string
}

type Bar struct {
	Name     string
	Greeting string
}

func newBar() *Bar {
	return &Bar{}
}

// PATH /foo
type FooController struct {
	at.RestController
	token jwt.Token

	MagicNumber   int64 `value:"${magic.number:888}"`
	DefaultNumber int64 `value:"${default.number:888}"`
}

func newFooController(token jwt.Token) *FooController {
	return &FooController{
		token: token,
	}
}

type ExampleController struct {
	at.RestController
}

type InvalidController struct{}

type FooBar struct {
	Name string
}

type FooBarService struct {
	fooBar *FooBar
}

func newFooBarService(fooBar *FooBar) *FooBarService {
	return &FooBarService{
		fooBar: fooBar,
	}
}

func (s *FooBarService) FooBar() *FooBar {
	return s.fooBar
}

func init() {
	log.SetLevel(log.DebugLevel)
	app.Register(&FooBar{Name: "fooBar"})
	app.Register(newFooBarService)
	app.Register(newHelloContextAware)
}

func (c *FooController) Before(ctx context.Context) {
	log.Debug("FooController.Before")

	ctx.Next()
}

func (c *FooController) PostLogin(request *UserRequest) (response model.Response, err error) {
	log.Debug("FooController.Login")

	// you make validate username and password first
	token, err := c.token.Generate(jwt.Map{
		"username": request.Username,
		"password": request.Password,
	}, 10, time.Minute)
	response = new(model.BaseResponse)
	response.SetData(token)

	return
}

// GetMagicNumber
func (c *FooController) GetMagic() (response model.Response, err error) {
	log.Debug("FooController.GetMagic")

	response = new(model.BaseResponse)
	data := map[string]interface{}{
		"magic_number":   c.MagicNumber,
		"default_number": c.DefaultNumber,
	}
	log.Debugf("default number: %v magic number: %v", c.DefaultNumber, c.MagicNumber)
	response.SetData(data)
	return
}

// POST /
func (c *FooController) Post(request *FooRequest) (response model.Response, err error) {
	log.Debug("FooController.Post")

	response = new(model.BaseResponse)
	if request.Name != "John" {
		return response, errors.New("only John is illegal name in this test")
	}
	response.SetData("Hello, " + request.Name)
	return
}

// GET /options/{options}
func (c *FooController) GetByOptions(options []string) (response model.Response) {
	type settings struct {
		Options []string
	}
	response = new(model.BaseResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("success")
	response.SetData(&settings{Options: options})
	return
}


// GET /foo/fs
func (c *FooController) GetFs(fs []*Foo) (response model.Response) {
	response = new(model.BaseResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("success")
	response.SetData(fs)
	return
}

// GET /name/{name}
func (c *FooController) GetByName(name string) (response model.Response) {
	type user struct {
		Name string
	}
	response = new(model.BaseResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("success")
	response.SetData(&user{Name: name})
	return
}

// GET /id/{id}
func (c *FooController) GetById(id int) string {
	log.Debugf("FooController.Get by id: %v", id)
	return "hello"
}

// GET /hello
func (c *FooController) GetHello(ctx context.Context) string {
	log.Debug(ctx.Annotations())
	ctx.SetURLParam("foo", "bar")
	ctx.SetURLParam("foo", "baz")
	log.Debug(ctx.URLParam("foo"))
	log.Debug("FooController.GetHello")
	return "hello"
}

// PUT /id/{id}/name/{name}/age/{age}
func (c *FooController) PutByIdNameAge(id int, name string, age int) error {
	log.Debugf("FooController.Put %v %v %v", id, name, age)
	return nil
}

// PATCH /id/{id}
func (c *FooController) PatchById(id int) error {
	log.Debug("FooController.Patch")
	return nil
}

// DELETE /id/{id}
func (c *FooController) DeleteById(id int) error {
	log.Debug("FooController.Delete ", id)
	return nil
}

func (c *FooController) GetBar(bar string) {
	log.Debugf("FooController.GetBar: %v", bar)
}

func (c *FooController) Options() {
	log.Debug("FooController.Options")
}

// GET /invalid
func (c *FooController) GetInvalid() (err error) {
	return fmt.Errorf("response invalid")
}

// GET /integer
func (c *FooController) GetInteger() (i int) {
	return 123
}

// GET /intPointer
func (c *FooController) GetIntPointer() (ip *int) {
	iv := 123
	ip = &iv
	return
}

// GET /intNullPointer
func (c *FooController) GetIntNilPointer() (ip *int) {
	return
}

func (c *FooController) GetError() error {
	return fmt.Errorf("buildconfigs.mio.io \"hello-world\" not found")
}

type FooRequestBody struct {
	at.RequestBody
	Name string `json:"name"`
}

func (c *FooController) GetRequestBody(request *FooRequestBody) string {
	return request.Name
}

type FooRequestForm struct {
	at.RequestForm
	Name string `json:"name"`
}

func (c *FooController) GetRequestForm(request *FooRequestForm) string {
	return request.Name
}

type EmbeddedParam struct {
	Name string `json:"name"`
}

// TODO: embedded query param does not work
type EmbeddedFooRequestParams struct {
	at.RequestParams
	EmbeddedParam
}

// TODO: embedded query param does not work
type FooRequestParams struct {
	at.RequestParams
	Name string `json:"name"`
}

func (c *FooController) GetRequestParams(request *FooRequestParams) string {
	return request.Name
}

type HelloContextAware struct {
	at.ContextAware

	context context.Context
}

func newHelloContextAware(context context.Context) *HelloContextAware {
	return &HelloContextAware{context: context}
}

func (c *FooController) GetContext(hca *HelloContextAware) string {
	return "testing context aware dependency injection"
}

func (c *FooController) GetErr() float32 {
	return 0.01
}


type Foo struct {
	Name string
}
// Get test get
func (c *FooController) GetInjection(foo *Foo) string {
	log.Debug(foo)
	return foo.Name
}

func (c *FooController) After() {
	log.Debug("FooController.After")
}

// BarController
type BarController struct {
	at.JwtRestController
}

func newBarController() *BarController {
	return &BarController{}
}

func (c *BarController) Get(request *BarRequest) (response model.Response, err error) {
	log.Debug("BarController.Get")
	response = new(model.BaseResponse)
	response.SetData("Hello, " + request.Name)

	return response, nil
}

type FoobarController struct {
	at.RestController
}

func newFoobarController() *FoobarController {
	return &FoobarController{}
}

func (c *FoobarController) Post(request *FoobarRequestForm) (response model.Response, err error) {
	response = new(model.BaseResponse)
	response.SetData("Hello, " + request.Name)

	return
}

func (c *FoobarController) Get(request *FoobarRequestParams) (response model.Response, err error) {
	response = new(model.BaseResponse)
	response.SetData("Hello, " + request.Name)

	return
}

// Define our controller, start with the name Foo, the first word of the Camelcase FooController is the controller name
// the lower cased foo will be the context mapping of the controller
// context mapping can be overwritten by FooController.ContextMapping
type HelloController struct {
	at.RestController
	fooBar *FooBar
}

func newHelloController(fooBar *FooBar) *HelloController {
	return &HelloController{
		fooBar: fooBar,
	}
}

// Get hello
// The first word of method is the http method GET, the rest is the context mapping hello
// in this method, the name Get means that the method context mapping is '/'
func (c *HelloController) Get() string {
	return "hello"
}

func (c *HelloController) GetWorld(ctx context.Context) {
	ctx.HTML("<h1>Hello World</h1>")
}

// Get /all
func (c *HelloController) GetAll(ctx context.Context) {

	data := []struct {
		Name string
		Age  int
	}{
		{
			Name: "John Doe",
			Age:  18,
		},
		{
			Name: "Zhang San",
			Age:  25,
		},
	}

	ctx.ResponseBody("success", data)
}

// GetMap
func (c *HelloController) GetMap() map[string]interface{} {
	response := make(map[string]interface{})
	response["message"] = "Hello from Map"
	return response
}

type oneTwoThreeController struct {
	at.RestController
}

func newOneTwoThreeController() *oneTwoThreeController {
	return &oneTwoThreeController{}
}

func (c *oneTwoThreeController) Get() string {
	return "from one-two-three controller"
}

func TestOneTwoThreeController(t *testing.T) {
	mu.Lock()
	testData := []struct {
		format string
		path   string
	}{
		{
			format: app.ContextPathFormatLowerCamel,
			path:   "/oneTwoThree",
		},
		{
			format: app.ContextPathFormatKebab,
			path:   "/one-two-three",
		},
		{
			format: app.ContextPathFormatSnake,
			path:   "/one_two_three",
		},
		{
			format: app.ContextPathFormatCamel,
			path:   "/OneTwoThree",
		},
	}
	for _, td := range testData {
		testApp := web.NewTestApp(newOneTwoThreeController).
			SetProperty(app.ContextPathFormat, td.format).
			SetProperty(logging.Level, log.InfoLevel).
			Run(t)
		t.Run("should response 200 when GET "+td.path, func(t *testing.T) {
			testApp.
				Get(td.path).
				Expect().Status(http.StatusOK)
		})
	}
	mu.Unlock()
}

// Define our controller, start with the name Foo, the first word of the Camelcase FooController is the controller name
// the lower cased foo will be the context mapping of the controller
// context path can be overwritten by the tag value of annotation at.ContextPath
type HelloViewController struct {
	at.RestController
	at.RequestMapping `value:"/"`

	fooBar *FooBar
}

func newHelloViewController(fooBar *FooBar) *HelloViewController {
	return &HelloViewController{
		fooBar: fooBar,
	}
}

// Get hello
// The first word of method is the http method GET, the rest is the context mapping hello
// in this method, the name Get means that the method context mapping is '/'
func (c *HelloViewController) Get(ctx context.Context) (err error) {
	err = ctx.View("index.html")
	return
}

// AnyTest
func (c *HelloViewController) AnyTest() string {
	return "Greeting from any method"
}

func TestWebViewApplicationWithProperties(t *testing.T) {
	mu.Lock()
	testApp := web.NewTestApp(newHelloViewController).
		SetProperty("web.view.enabled", true).
		Run(t)
	t.Run("should response 200 when GET /", func(t *testing.T) {
		testApp.
			Get("/").
			Expect().Status(http.StatusOK)
	})
	mu.Unlock()
}

func TestWebViewApplicationWithArgs(t *testing.T) {
	mu.Lock()
	testApp := web.NewTestApp(newHelloViewController).
		SetProperty("server.port", 8080).
		SetProperty("web.view.enabled", true).
		Run(t)
	t.Run("should response 200 when GET /", func(t *testing.T) {
		testApp.
			Get("/").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response 200 when GET /test", func(t *testing.T) {
		testApp.
			Get("/test").
			Expect().Status(http.StatusOK)
	})
	mu.Unlock()
}

func TestApplicationWithoutController(t *testing.T) {
	mu.Lock()
	testApp := web.NewTestApp().
		Run(t)

	t.Run("should response 404 when GET /", func(t *testing.T) {
		testApp.
			Get("/").
			Expect().Status(http.StatusNotFound)
	})
	mu.Unlock()
}

type circularFoo struct {
	circularBar *circularBar
}

func newCircularFoo(circularBar *circularBar) *circularFoo {
	return &circularFoo{
		circularBar: circularBar,
	}
}

type circularBar struct {
	circularFoo *circularFoo
}

func newCircularBar(circularFoo *circularFoo) *circularBar {
	return &circularBar{
		circularFoo: circularFoo,
	}
}

type circularDiController struct {
	circularFoo *circularFoo
}

func newCircularDiController(circularFoo *circularFoo) *circularDiController {
	return &circularDiController{
		circularFoo: circularFoo,
	}
}

func TestApplicationWithCircularDI(t *testing.T) {
	mu.Lock()
	testApp := web.NewTestApp(newCircularDiController).
		Run(t)

	t.Run("should response 404 when GET /", func(t *testing.T) {
		testApp.
			Get("/").
			Expect().Status(http.StatusNotFound)
	})
	mu.Unlock()
}

func TestWebApplication(t *testing.T) {
	mu.Lock()
	foo := &Foo{Name: "test injection"}
	app.Register(foo)
	testApp := web.NewTestApp(newHelloController, newFooController, newBarController, newFoobarController).
		SetProperty(app.ProfilesInclude, jwt.Profile, locale.Profile, logging.Profile).
		Run(t)

	t.Run("should response 200 when GET /hello/all", func(t *testing.T) {
		testApp.
			Request(http.MethodGet, "/hello/all").
			Expect().Status(http.StatusOK).
			Body().Contains("Success").Contains("John Doe").Contains("Zhang San")
	})

	t.Run("should response 200 when GET /hello", func(t *testing.T) {
		testApp.
			Get("/hello/").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response 200 when GET /map", func(t *testing.T) {
		testApp.
			Get("/hello/map").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response 200 when GET /hello/world", func(t *testing.T) {
		testApp.
			Get("/hello/world").
			Expect().Status(http.StatusOK)
	})

	t.Run("should response 你好, 世界 with language zh-CN", func(t *testing.T) {
		// test cn-ZH
		testApp.Get("/hello").
			WithHeader("Accept-Language", "zh-CN").
			Expect().Status(http.StatusOK).
			Body().Contains("你好, 世界")
	})

	t.Run("should response Success with language en-US", func(t *testing.T) {
		// test en-US
		testApp.Get("/hello").
			WithHeader("Accept-Language", "en-US").
			Expect().
			Status(http.StatusOK).
			Body().Contains("Hello, World")
	})

	//t.Run("should pass health check", func(t *testing.T) {
	//	app.Get("/health").Expect().Status(http.StatusOK)
	//})
	t.Run("should inject into controller method", func(t *testing.T) {
		testApp.Get("/foo/injection").
			Expect().Status(http.StatusOK).Body().Contains(foo.Name)
	})

	t.Run("should login with username and password", func(t *testing.T) {
		testApp.Post("/foo/login").
			WithJSON(&UserRequest{Username: "johndoe", Password: "iHop91#15"}).
			Expect().Status(http.StatusOK).
			Body().NotEqual("")
	})

	t.Run("should failed to pass validation", func(t *testing.T) {
		testApp.Post("/foo/login").
			WithJSON(&UserRequest{Username: "johndoe"}).
			Expect().Status(http.StatusBadRequest)
	})

	t.Run("should return success after POST /foo", func(t *testing.T) {
		testApp.Post("/foo").
			WithJSON(&FooRequest{Name: "John"}).
			Expect().Status(http.StatusOK).
			Body().Contains("Hello")
	})

	t.Run("should return StatusInternalServerError after POST /foo with illegal name", func(t *testing.T) {
		testApp.Post("/foo").
			WithJSON(&FooRequest{Name: "Illegal name"}).
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return success after GET /foo/hello", func(t *testing.T) {
		testApp.Get("/foo/hello").
			WithJSON(&FooRequest{Name: "John"}).
			Expect().Status(http.StatusOK).
			Body().Equal("Hello, World")
	})

	t.Run("should parse request body GET /foo/requestBody", func(t *testing.T) {
		testApp.Get("/foo/requestBody").
			WithJSON(&FooRequestBody{Name: "foo"}).
			Expect().Status(http.StatusOK).
			Body().Equal("foo")
	})

	t.Run("should parse request body GET /foo/context", func(t *testing.T) {
		testApp.Get("/foo/context").
			Expect().Status(http.StatusOK)
	})

	t.Run("should parse request body GET /foo/err", func(t *testing.T) {
		testApp.Get("/foo/err").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get magic number GET /foo/magic", func(t *testing.T) {
		testApp.Get("/foo/magic").
			Expect().Status(http.StatusOK)
	})

	t.Run("should parse request body GET /foo/requestForm", func(t *testing.T) {
		testApp.Get("/foo/requestForm").
			WithFormField("name", "foo").
			Expect().Status(http.StatusOK).
			Body().Equal("")
	})

	t.Run("should parse request body GET /foo/requestParams", func(t *testing.T) {
		testApp.Get("/foo/requestParams").
			WithQuery("name", "foo").
			Expect().Status(http.StatusOK).
			Body().Equal("foo")
	})

	t.Run("should return http.StatusUnauthorized after GET /bar", func(t *testing.T) {
		testApp.Get("/bar").
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should return http.StatusInternalServerError when input form field validation failed", func(t *testing.T) {
		// test request form
		testApp.Post("/foobar").
			WithFormField("name", "John Doe").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return (http.StatusOK on /foobar", func(t *testing.T) {
		//  test request query
		testApp.Get("/foobar").
			WithQuery("name", "John Doe").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /bar with jwt token", func(t *testing.T) {
		log.Println(io.GetWorkDir())
		jwtToken := jwt.NewJwtToken(&jwt.Properties{
			PrivateKeyPath: "config/ssl/app.rsa",
			PublicKeyPath:  "config/ssl/app.rsa.pub",
		})
		// test jwt
		token, err := jwtToken.Generate(jwt.Map{
			"username": "johndoe",
			"password": "PA$$W0RD",
		}, 100, time.Millisecond)
		if err == nil {

			t := fmt.Sprintf("Bearer %v", token)

			testApp.Get("/bar").
				WithHeader("Authorization", t).
				Expect().Status(http.StatusOK)

			time.Sleep(2 * time.Second)

			testApp.Get("/bar").
				WithHeader("Authorization", t).
				Expect().Status(http.StatusUnauthorized)
		}
	})

	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		testApp.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", 123456).
			WithPath("name", "Mike").
			WithPath("age", 18).
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		testApp.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", 0).
			WithPath("name", "").
			WithPath("age", 0).
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		testApp.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", " ").
			WithPath("name", " ").
			WithPath("age", " ").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		testApp.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", "").
			WithPath("name", "").
			WithPath("age", " ").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusOK on /foo with PUT, PATCH, DELETE methods", func(t *testing.T) {
		testApp.Put("/foo/id/{id}/name/{name}/age/{age}").
			WithPath("id", "").
			WithPath("name", "").
			WithPath("age", "").
			Expect().Status(http.StatusTemporaryRedirect)
	})

	t.Run("should Get foo by id", func(t *testing.T) {
		testApp.Get("/foo/id/{id}").
			WithPath("id", 123).
			Expect().Status(http.StatusOK)
	})

	t.Run("should Get by options", func(t *testing.T) {
		testApp.Get("/foo/fs").
			WithJSON([]*Foo{
				{
					Name: "foo",
				},
				{
					Name: "bar",
				},
		}).
			Expect().Status(http.StatusOK).
			Body().Contains("foo")
	})

	t.Run("should Get by options", func(t *testing.T) {
		testApp.Get("/foo/options/{options}").
			WithPath("options", "mars,earth,mercury,jupiter").
			Expect().Status(http.StatusOK).
			Body().Contains("mercury")
	})

	t.Run("should Get foo by name", func(t *testing.T) {
		testApp.Get("/foo/name/{name}").
			WithPath("name", "Peter Phil").
			Expect().Status(http.StatusOK).
			Body().Contains("Peter Phil")
	})

	t.Run("should Get foo by name", func(t *testing.T) {
		testApp.Get("/foo/name/{name}").
			WithPath("name", "张三").
			Expect().Status(http.StatusOK).
			Body().Contains("张三")
	})

	t.Run("should Patch foo by id", func(t *testing.T) {
		testApp.Patch("/foo/id/{id}").
			WithPath("id", 456).
			Expect().Status(http.StatusOK)
	})

	t.Run("should Delete foo by id", func(t *testing.T) {
		testApp.Delete("/foo/id/{id}").
			WithPath("id", 789).
			Expect().Status(http.StatusOK)
	})

	t.Run("should call Options method", func(t *testing.T) {
		testApp.Options("/foo").
			Expect().Status(http.StatusOK)
	})

	// TODO: this test case is deleted due to the implementation of the anonymous struct that accept the route annotation.
	//t.Run("should return failure for unsupported args", func(t *testing.T) {
	//	testApp.Get("/foo/bar").
	//		Expect().Status(http.StatusInternalServerError)
	//})

	t.Run("should return invalid response", func(t *testing.T) {
		testApp.Get("/foo/invalid").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should return integer", func(t *testing.T) {
		testApp.Get("/foo/integer").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return integer pointer", func(t *testing.T) {
		testApp.Get("/foo/intPointer").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return integer nil pointer", func(t *testing.T) {
		testApp.Get("/foo/intNilPointer").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return error message", func(t *testing.T) {
		testApp.Get("/foo/error").
			Expect().Status(http.StatusInternalServerError)
	})
	mu.Unlock()
}

//func TestInvalidController(t *testing.T)  {
//	ta := new(testApplication)
//	err := ta.Init(new(InvalidController))
//	err, ok := err.(*system.InvalidControllerError)
//	assert.Equal(t, ok, true)
//}

func TestNewApplication(t *testing.T) {
	mu.Lock()
	app.Register(new(ExampleController))

	testApp := web.NewApplication().
		SetAddCommandLineProperties(true).
		SetProperty("server.port", 8080).
		SetProperty(app.BannerDisabled, true)
	t.Run("should init web application", func(t *testing.T) {
		assert.NotEqual(t, nil, testApp)
	})

	t.Run("should get interface name", func(t *testing.T) {
		type myInterface interface {
		}

		mi := new(myInterface)

		typ, err := reflector.GetType(mi)
		assert.Equal(t, nil, err)
		assert.Equal(t, "myInterface", typ.Name())
	})

	go testApp.Run()
	time.Sleep(time.Second)
	mu.Unlock()
}

func TestAnonymousController(t *testing.T) {
	mu.Lock()
	t.Run("should failed to register anonymous controller", func(t *testing.T) {
		testApp := web.RunTestApplication(t, (*Bar)(nil))
		assert.NotEqual(t, nil, testApp)
	})

	t.Run("should failed to register anonymous controller", func(t *testing.T) {
		testApp := web.RunTestApplication(t, newBar)
		assert.NotEqual(t, nil, testApp)
	})
	mu.Unlock()
}

type fakeConditionalJwtMiddleware struct {
	at.Middleware
	at.UseJwt
}

func newConditionalFakeJwtMiddleware() *fakeConditionalJwtMiddleware {
	return &fakeConditionalJwtMiddleware{}
}

// CheckJwt
func (m *fakeConditionalJwtMiddleware) CheckJwt(at struct{ at.MiddlewareHandler
}, ctx context.Context)  {
	log.Debug("fakeConditionalJwtMiddleware.CheckJwt()")
	if ctx.URLParam("token") == "" {
		ctx.StatusCode(http.StatusUnauthorized)
		return
	}
	ctx.Next()
	return
}

type fakeMethodConditionalJwtMiddleware struct {
	at.Middleware
}

func newMethodConditionalFakeJwtMiddleware() *fakeMethodConditionalJwtMiddleware {
	return &fakeMethodConditionalJwtMiddleware{}
}

// CheckJwt
func (m *fakeMethodConditionalJwtMiddleware) CheckJwt(at struct{
	at.MiddlewareHandler
	at.UseJwt
}, ctx context.Context)  {
	log.Debug("fakeMethodConditionalJwtMiddleware.CheckJwt()")
	if ctx.URLParam("token") == "" {
		ctx.StatusCode(http.StatusUnauthorized)
		return
	}
	ctx.Next()
	return
}


type fooMiddleware struct {
	at.Middleware
}

func newFooMiddleware() *fooMiddleware {
	return &fooMiddleware{}
}

// Logging is the middleware handler,it support dependency injection, method annotation
// middleware handler can be annotated to specific purpose or general purpose
func (m *fooMiddleware) Logging( _ struct{at.MiddlewareHandler `value:"/" `}, ctx context.Context) {

	log.Infof("[logging middleware] %v", ctx.GetCurrentRoute())

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}
// Logging is the middleware handler,it support dependency injection, method annotation
// middleware handler can be annotated to specific purpose or general purpose
func (m *fooMiddleware) PostLogging( _ struct{at.MiddlewarePostHandler `value:"/" `}, ctx context.Context) {

	log.Infof("[logging middleware] %v", ctx.GetCurrentRoute())

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

type CustomResponse struct {
	at.ResponseBody
	model.BaseResponseInfo
	Data interface{} `json:"data"`
}

type customRouterController struct {
	at.RestController

	at.RequestMapping `value:"/custom"`
}

func newCustomRouterController() *customRouterController {
	return new(customRouterController)
}

func (c *customRouterController) PathVariable(
	at struct {
	// at.GetMapping is an annotation to define request mapping for http method GET /{id}/and/{name}
	at.GetMapping `value:"/{id}/name/{name}"`
}, id int, name string) (response *CustomResponse, err error) {
	response = new(CustomResponse)
	log.Infof("PathParamIdAndName: %v", at.AtValue)
	switch id {
	case 0:
		response.Code = http.StatusNotFound
		err = fmt.Errorf("not found")
	case 1:
		err = fmt.Errorf("wrong id")
	case 2:
	default:
		response.Code = http.StatusOK
		response.Message = "Success"
		response.Data = fmt.Sprintf("https://hidevops.io/%v/%v", id, name)
	}

	return
}

func (c *customRouterController) Delete(
	at struct {
	// at.GetMapping is an annotation to define request mapping for http method GET /{id}/and/{name}
	at.DeleteMapping `value:"/{id}"`
	// uses jwt auth, it can be placed in struct of method
	at.UseJwt
}, id int,) (response *CustomResponse, err error) {
	response = new(CustomResponse)
	log.Infof("Delete: %v", at.DeleteMapping.AtValue)

	response.Code = http.StatusOK
	response.Message = "Success"
	return
}

// BeforeMethod
func (c *customRouterController) BeforeMethod(at struct{ at.BeforeMethod}, ctx context.Context)  {
	ctx.Next()
	return
}

// AfterMethod
func (c *customRouterController) AfterMethod(at struct{ at.AfterMethod })  {
	return
}

type jwtAuthTestController struct {
	at.RestController
	at.UseJwt
	at.RequestMapping `value:"jwt-auth" `
}

func newJwtAuthTestController() *jwtAuthTestController {
	return &jwtAuthTestController{}
}

// Get
func (c *jwtAuthTestController) Get(at struct{ at.GetMapping `value:"/"` }) string {
	return "Get from jwt auth test controller"
}

// Delete
func (c *jwtAuthTestController) Delete(at struct{ at.DeleteMapping `value:"/"` }) string  {
	return "Delete from jwt auth test controller"
}

type regularTestController struct {
	at.RestController
	at.RequestMapping `value:"regular" `
}

func init() {
	app.Register(newRegularTestController)
}

func newRegularTestController() *regularTestController {
	return &regularTestController{}
}

// Get
func (c *regularTestController) Get(at struct{ at.GetMapping `value:"/"` }) string  {
	return "regularTestController.Get()"
}

func TestCustomRouter(t *testing.T) {
	mu.Lock()
	app.Register(newFooMiddleware)
	testApp := web.NewTestApp(newCustomRouterController).
		SetProperty("server.context_path", "/test").
		Run(t)
	testApp.Get("/test/custom/0/name/hiboot").Expect().Status(http.StatusNotFound)
	testApp.Get("/test/custom/123/name/hiboot").Expect().Status(http.StatusOK)
	testApp.Get("/test/custom/1/name/hiboot").Expect().Status(http.StatusInternalServerError)
	testApp.Get("/test/custom/2/name/hiboot").Expect().Status(http.StatusOK)
	mu.Unlock()
}

// test HttpMethodSubscriber
type fakeSubscriber struct {
	at.HttpMethodSubscriber
}

func newFakeSubscriber() *fakeSubscriber {
	return &fakeSubscriber{}
}

func (s *fakeSubscriber) Subscribe(atc *annotation.Annotations, atm *annotation.Annotations)  {
	log.Debug("subscribe")
}

func TestMiddlewareAnnotation(t *testing.T) {
	mu.Lock()
	app.Register(
		newCustomRouterController,
		newFooMiddleware,
		newConditionalFakeJwtMiddleware,
		newJwtAuthTestController,
		newRegularTestController,
		newMethodConditionalFakeJwtMiddleware,
		newFakeSubscriber)
	testApp := web.NewTestApp(newCustomRouterController).
		Run(t)

	t.Run("should get resource from custom controller", func(t *testing.T) {
		testApp.Get("/custom/123/name/hiboot").
			Expect().Status(http.StatusOK)
	})

	t.Run("should delete resource from custom controller with jwt token", func(t *testing.T) {
		testApp.Delete("/custom/123").
			WithQuery("token", "fake-token").
			Expect().Status(http.StatusOK)
	})

	t.Run("should not delete resource from custom controller without jwt token", func(t *testing.T) {
		testApp.Delete("/custom/123").
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should get resource from regular controller", func(t *testing.T) {
		testApp.Get("/regular").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get resource from jwt-auth controller with jwt token", func(t *testing.T) {
		testApp.Get("/jwt-auth").
			WithQuery("token", "fake-token").
			Expect().Status(http.StatusOK)
	})

	t.Run("should delete resource from jwt-auth controller with jwt token", func(t *testing.T) {
		testApp.Delete("/jwt-auth").
			WithQuery("token", "fake-token").
			Expect().Status(http.StatusOK)
	})

	t.Run("should not delete resource from jwt-auth controller without jwt token", func(t *testing.T) {
		testApp.Delete("/jwt-auth").
			Expect().Status(http.StatusUnauthorized)
	})
	mu.Unlock()
}

// --- test file server
type publicTestController struct {
	at.RestController
	at.RequestMapping `value:"/public"`
}

func newPublicTestController() *publicTestController {
	return &publicTestController{}
}

// UI serve static resource via context StaticResource method
func (c *publicTestController) UI(at struct{ at.GetMapping `value:"/ui/*"`; at.FileServer `value:"/ui"` }, ctx context.Context) {
	return
}

// UI serve static resource via context StaticResource method
func (c *publicTestController) UIIndex(at struct{ at.GetMapping `value:"/ui"`; at.FileServer `value:"/ui"` }, ctx context.Context) {
	return
}

// UI serve static resource via context StaticResource method
func (c *publicTestController) NotFound(at struct{ at.GetMapping `value:"/foo"`; at.FileServer `value:"/ui"` }, ctx context.Context) {
	ctx.WrapHandler( http.NotFoundHandler() )
}

func TestController(t *testing.T) {
	mu.Lock()
	testApp := web.NewTestApp(t, newPublicTestController).
		SetProperty("server.port", "8085").Run(t)

	t.Run("should get index.html ", func(t *testing.T) {
		testApp.Get("/public/ui").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get hello.txt ", func(t *testing.T) {
		testApp.Get("/public/ui/hello.txt").
			Expect().Status(http.StatusOK)
	})

	t.Run("should report not found ", func(t *testing.T) {
		testApp.Get("/public/foo").
			Expect().Status(http.StatusNotFound)
	})
	mu.Unlock()
}
