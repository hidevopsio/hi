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

// Package factory provides InstantiateFactory and ConfigurableFactory interface
package factory

import (
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/system/scheduler"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"reflect"
)

const (
	// InstantiateFactoryName is the instance name of factory.instantiateFactory
	InstantiateFactoryName = "factory.instantiateFactory"
	// ConfigurableFactoryName is the instance name of factory.configurableFactory
	ConfigurableFactoryName = "factory.configurableFactory"
)

// Factory interface
type Factory interface{}

type Instance interface {
	Get(params ...interface{}) (retVal interface{})
	Set(params ...interface{}) (err error)
	Items() map[string]interface{}
}

// InstantiateFactory instantiate factory interface
type InstantiateFactory interface {
	Initialized() bool
	SetInstance(params ...interface{}) (err error)
	GetInstance(params ...interface{}) (retVal interface{})
	GetInstances(params ...interface{}) (retVal []*MetaData)
	Items() map[string]interface{}
	Append(i ...interface{})
	AppendComponent(c ...interface{})
	BuildComponents() (err error)
	Builder() (builder system.Builder)
	GetProperty(name string) interface{}
	SetProperty(name string, value interface{}) InstantiateFactory
	SetDefaultProperty(name string, value interface{}) InstantiateFactory
	DefaultProperties() map[string]interface{}
	InjectIntoFunc(object interface{}) (retVal interface{}, err error)
	InjectIntoMethod(owner, object interface{}) (retVal interface{}, err error)
	InjectDefaultValue(object interface{}) error
	InjectIntoObject(object interface{}) error
	InjectDependency(object interface{}) (err error)
	Replace(name string) interface{}
	InjectContextAwareObjects(ctx context.Context, dps []*MetaData) (runtimeInstance Instance, err error)
}

// ConfigurableFactory configurable factory interface
type ConfigurableFactory interface {
	InstantiateFactory
	SystemConfiguration() *system.Configuration
	Configuration(name string) interface{}
	BuildProperties() (systemConfig *system.Configuration, err error)
	StartSchedulers(schedulerServices []*MetaData) (schedulers []*scheduler.Scheduler)
	Build(configs []*MetaData)
}

// Configuration configuration interface
type Configuration interface {
}

type depsMap map[string][]string

// Deps the dependency mapping of configuration
type Deps struct {
	deps depsMap
}

func (c *Deps) ensure() {
	if c.deps == nil {
		c.deps = make(depsMap)
	}
}

// Get get the dependencies mapping
func (c *Deps) Get(name string) (deps []string) {
	c.ensure()

	deps = c.deps[name]

	return
}

// Set set dependencies
func (c *Deps) Set(dep interface{}, value []string) {
	c.ensure()
	var name string
	val := reflect.ValueOf(dep)
	kind := val.Kind()
	switch kind {
	case reflect.Func:
		name = reflector.GetFuncName(dep)
	case reflect.String:
		name = dep.(string)
	default:
		return
	}
	c.deps[name] = value

	return
}

// CastMetaData cast object to *factory.MetaData
func CastMetaData(object interface{}) (metaData *MetaData) {
	switch object.(type) {
	case *MetaData:
		metaData = object.(*MetaData)
	}
	return
}
