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

// Package app provides abstract layer for cli/web application
package app

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/system/scheduler"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
)

const (
	// ApplicationContextName is the application context instance name
	ApplicationContextName      = "app.applicationContext"
	ContextPathFormat           = "server.context_path_format"
	ContextPathFormatKebab      = "kebab"
	ContextPathFormatSnake      = "snake"
	ContextPathFormatCamel      = "camel"
	ContextPathFormatLowerCamel = "lower-camel"
)

// Application is the base application interface
type Application interface {
	Initialize() error
	// TODO: remove it from factory as system.build can be injected directly
	SetProperty(name string, value ...interface{}) Application
	GetProperty(name string) (value interface{}, ok bool)
	SetAddCommandLineProperties(enabled bool) Application
	Run()
}

// ApplicationContext is the alias interface of Application
type ApplicationContext interface {
	RegisterController(controller interface{}) error
	Use(handlers ...context.Handler)
	GetProperty(name string) (value interface{}, ok bool)
	GetInstance(params ...interface{}) (instance interface{})
}

// BaseApplication is the base application
type BaseApplication struct {
	WorkDir             string
	configurations      cmap.ConcurrentMap
	instances           cmap.ConcurrentMap
	potatoes            cmap.ConcurrentMap
	configurableFactory factory.ConfigurableFactory
	systemConfig        *system.Configuration
	postProcessor       *postProcessor
	defaultProperties   cmap.ConcurrentMap
	mu                  sync.Mutex
	// SetAddCommandLineProperties
	addCommandLineProperties bool

	schedulers []*scheduler.Scheduler
}

var (
	configContainer    []*factory.MetaData
	componentContainer []*factory.MetaData
	// Profiles include profiles initially
	Profiles           []string

	// ErrInvalidObjectType indicates that configuration type is invalid
	ErrInvalidObjectType = errors.New("[app] invalid Configuration type, one of app.Configuration need to be embedded")

	banner = `
______  ____________             _____
___  / / /__(_)__  /_______________  /_
__  /_/ /__  /__  __ \  __ \  __ \  __/   
_  __  / _  / _  /_/ / /_/ / /_/ / /_     Hiboot Application Framework
/_/ /_/  /_/  /_.___/\____/\____/\__/     https://github.com/hidevopsio/hiboot

`
)

// PrintStartupMessages prints startup messages
func (a *BaseApplication) PrintStartupMessages() {
	if !a.systemConfig.App.Banner.Disabled {
		if a.systemConfig.App.Banner.Custom != "" {
			fmt.Print(a.systemConfig.App.Banner.Custom)
		} else {
			fmt.Print(banner)
		}
	}
}

// SetProperty set application property
// should be able to set property from source code by SetProperty, it can be override by program argument, e.g. myapp --app.profiles.active=dev
func (a *BaseApplication) SetProperty(name string, value ...interface{}) Application {
	var val interface{}
	if len(value) == 1 {
		val = value[0]
	} else {
		val = value
	}

	kind := reflect.TypeOf(val).Kind()
	if kind == reflect.String && strings.Contains(val.(string), ",") {
		val = strings.SplitN(val.(string), ",", -1)
	}
	a.defaultProperties.Set(name, val)

	return a
}

// GetProperty get application property
func (a *BaseApplication) GetProperty(name string) (value interface{}, ok bool) {
	value, ok = a.defaultProperties.Get(name)
	return
}

// Initialize init application
func (a *BaseApplication) Initialize() (err error) {
	log.SetLevel(log.InfoLevel)
	a.defaultProperties = cmap.New()
	a.configurations = cmap.New()
	a.instances = cmap.New()
	// set add command line properties to true as default
	a.SetAddCommandLineProperties(true)
	return nil
}

// Initialize init application
func (a *BaseApplication) Build() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.WorkDir = io.GetWorkDir()

	instantiateFactory := instantiate.NewInstantiateFactory(a.instances, componentContainer, a.defaultProperties)
	configurableFactory := autoconfigure.NewConfigurableFactory(instantiateFactory, a.configurations)
	a.configurableFactory = configurableFactory

	a.postProcessor = newPostProcessor(instantiateFactory)
	a.systemConfig, _ = configurableFactory.BuildProperties()

	// set logging level
	log.SetLevel(a.systemConfig.Logging.Level)
}

// SystemConfig returns application config
func (a *BaseApplication) SystemConfig() *system.Configuration {
	return a.systemConfig
}

// BuildConfigurations get BuildConfigurations
func (a *BaseApplication) BuildConfigurations() (err error) {
	// build configurations
	a.configurableFactory.Build(configContainer)
	// build components
	err = a.configurableFactory.BuildComponents()

	// Start Scheduler after build
	schedulerServices := a.configurableFactory.GetInstances(at.EnableScheduling{})
	a.schedulers = a.configurableFactory.StartSchedulers(schedulerServices)

	return
}

// ConfigurableFactory get ConfigurableFactory
func (a *BaseApplication) ConfigurableFactory() factory.ConfigurableFactory {
	return a.configurableFactory
}

// AfterInitialization post initialization
func (a *BaseApplication) AfterInitialization(configs ...cmap.ConcurrentMap) {
	// pass user's instances
	a.postProcessor.Init()
	a.postProcessor.AfterInitialization()
	log.Infof("command line properties is enabled: %t", a.addCommandLineProperties)
}

// RegisterController register controller by interface
func (a *BaseApplication) RegisterController(controller interface{}) error {
	return nil
}

// Use use middleware handlers
func (a *BaseApplication) Use(handlers ...context.Handler) {
}

// SetAddCommandLineProperties set add command line properties to be enabled or disabled
func (a *BaseApplication) SetAddCommandLineProperties(enabled bool) Application {
	a.addCommandLineProperties = enabled
	return a
}

// Run run the application
func (a *BaseApplication) Run() {
	log.Warn("application is not implemented!")
}

// GetInstance get application instance by name
func (a *BaseApplication) GetInstance(params ...interface{}) (instance interface{}) {
	if a.configurableFactory != nil {
		instance = a.configurableFactory.GetInstance(params...)
	}
	return
}
