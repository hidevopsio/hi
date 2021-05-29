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
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/kataras/iris"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
)

const (
	pathSep = "/"

	beforeMethod       = "Before"
	afterMethod        = "After"
	applicationContext = "app.applicationContext"
)

type webApp struct {
	*iris.Application
}

func newWebApplication() *webApp {
	return &webApp{
		Application: iris.New(),
	}
}

// Application is the struct of web Application
type application struct {
	app.BaseApplication
	webApp     *webApp
	jwtEnabled bool
	//anonControllers []interface{}
	//jwtControllers  []interface{}
	controllers []interface{}
	dispatcher  *Dispatcher
	//controllerMap map[string][]interface{}
	startUpTime time.Time
}

var (
	// controllers global controllers container
	registeredControllers []interface{}
	compiledRegExp        = regexp.MustCompile(`\{(.*?)\}`)

	// ErrControllersNotFound controller not found
	ErrControllersNotFound = errors.New("[app] controllers not found")

	// ErrInvalidController invalid controller
	ErrInvalidController = errors.New("[app] invalid controller")
)

// SetProperty set application property
func (a *application) SetProperty(name string, value ...interface{}) app.Application {
	a.BaseApplication.SetProperty(name, value...)
	return a
}

// Initialize init application
func (a *application) Initialize() error {
	return a.BaseApplication.Initialize()
}

// Run run web application
func (a *application) Run() {
	serverPort := ":8080"
	// build HiBoot Application
	err := a.build()
	conf := a.SystemConfig()
	if conf != nil && err == nil {
		if conf.Server.Port != "" {
			serverPort = fmt.Sprintf(":%v", conf.Server.Port)
		}
		log.Infof("Hiboot started on port(s) http://localhost%v", serverPort)
		timeDiff := time.Since(a.startUpTime)
		log.Infof("Started %v in %f seconds", conf.App.Name, timeDiff.Seconds())
		// build web app
		a.webApp.Configure(iris.WithConfiguration(defaultConfiguration()))
		err = a.webApp.Build()

		// handler to Serve HTTP
		http.Handle("/", a.webApp)

		// serve web app with server port, default port number is 8080
		if err == nil {
			err = http.ListenAndServe(serverPort, nil)
			log.Debug(err)
		}
	}
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Init init web application
func (a *application) build() (err error) {

	a.Build()

	// set custom properties
	a.PrintStartupMessages()

	systemConfig := a.SystemConfig()
	// should do deduplication
	systemConfig.App.Profiles.Include = append(systemConfig.App.Profiles.Include, app.Profiles...)
	systemConfig.App.Profiles.Include = unique(systemConfig.App.Profiles.Include)
	if systemConfig != nil {
		log.Infof("Starting Hiboot web application %v version %v on localhost with PID %v", systemConfig.App.Name, systemConfig.App.Version, os.Getpid())
		log.Infof("Working directory: %v", a.WorkDir)
		log.Infof("The following profiles are active: %v, %v", systemConfig.App.Profiles.Active, systemConfig.App.Profiles.Include)
	}
	log.Infof("Initializing Hiboot Application")
	f := a.ConfigurableFactory()
	f.AppendComponent(app.ApplicationContextName, a)
	f.SetInstance(app.ApplicationContextName, a)

	// fill controllers into component container
	for _, ctrl := range a.controllers {
		f.AppendComponent(ctrl)
	}

	// build auto configurations
	err = a.BuildConfigurations()
	if err != nil {
		return
	}

	// create dispatcher
	a.dispatcher = a.GetInstance(Dispatcher{}).(*Dispatcher)

	// first register anon controllers
	err = a.RegisterController(at.RestController{})
	if err == nil {
		// call AfterInitialization with factory interface
		a.AfterInitialization()
	}
	return
}

// RegisterController register controller, e.g. at.RestController, jwt.Controller, or other customized controller
func (a *application) RegisterController(controller interface{}) error {
	middleware := a.ConfigurableFactory().GetInstances(at.Middleware{})
	log.Debug(middleware)
	// get from controller map
	// parse controller type
	controllers := a.ConfigurableFactory().GetInstances(controller)
	if controllers != nil {
		return a.dispatcher.register(controllers, middleware)
	}
	return ErrControllersNotFound
}

// Use apply middleware
func (a *application) Use(handlers ...context.Handler) {
	// pass user's instances
	for _, hdl := range handlers {
		a.webApp.Use(Handler(hdl))
	}
}

func (a *application) initialize(controllers ...interface{}) (err error) {
	io.EnsureWorkDir(3, "config/application.yml")

	// new iris app
	a.webApp = newWebApplication()
	app.Register(a.webApp)

	err = a.Initialize()

	if err == nil {
		if len(controllers) == 0 {
			a.controllers = registeredControllers
		} else {
			a.controllers = controllers
		}
	}
	return
}

// SetAddCommandLineProperties set add command line properties to be enabled or disabled
func (a *application) SetAddCommandLineProperties(enabled bool) app.Application {
	a.BaseApplication.SetAddCommandLineProperties(enabled)
	return a
}

// RestController register rest controller to controllers container
// Deprecated: please use app.Register() instead
var RestController = app.Register

// NewApplication create new web application instance and init it
func NewApplication(controllers ...interface{}) app.Application {
	log.SetLevel("error") // set debug level to error first
	a := new(application)
	app.Register(a)
	a.startUpTime = time.Now()
	_ = a.initialize(controllers...)
	return a
}
