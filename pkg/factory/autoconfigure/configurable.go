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

// Package autoconfigure implement ConfigurableFactory
package autoconfigure

import (
	"errors"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/system/scheduler"
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"os"
	"reflect"
	"strings"
)

const (
	// System configuration name
	System = "system"

	// PropAppProfilesActive is the property name "app.profiles.active"
	PropAppProfilesActive = "app.profiles.active"

	// EnvAppProfilesActive is the environment variable name APP_PROFILES_ACTIVE
	EnvAppProfilesActive = "APP_PROFILES_ACTIVE"

	// PostfixConfiguration is the Configuration postfix
	PostfixConfiguration = "Configuration"

	defaultProfileName = "default"
)

var (
	// ErrInvalidMethod method is invalid
	ErrInvalidMethod = errors.New("[factory] method is invalid")

	// ErrFactoryCannotBeNil means that the InstantiateFactory can not be nil
	ErrFactoryCannotBeNil = errors.New("[factory] InstantiateFactory can not be nil")

	// ErrFactoryIsNotInitialized means that the InstantiateFactory is not initialized
	ErrFactoryIsNotInitialized = errors.New("[factory] InstantiateFactory is not initialized")

	// ErrInvalidObjectType means that the Configuration type is invalid, it should embeds app.Configuration
	ErrInvalidObjectType = errors.New("[factory] invalid Configuration type, one of app.Configuration need to be embedded")

	// ErrConfigurationNameIsTaken means that the configuration name is already taken
	ErrConfigurationNameIsTaken = errors.New("[factory] configuration name is already taken")

	// ErrComponentNameIsTaken means that the component name is already taken
	ErrComponentNameIsTaken = errors.New("[factory] component name is already taken")
)

type configurableFactory struct {
	at.Qualifier `value:"factory.configurableFactory"`

	factory.InstantiateFactory
	configurations cmap.ConcurrentMap
	systemConfig   *system.Configuration

	preConfigureContainer  []*factory.MetaData
	configureContainer     []*factory.MetaData
	postConfigureContainer []*factory.MetaData
	builder                system.Builder
}

// NewConfigurableFactory is the constructor of configurableFactory
func NewConfigurableFactory(instantiateFactory factory.InstantiateFactory, configurations cmap.ConcurrentMap) factory.ConfigurableFactory {
	f := &configurableFactory{
		InstantiateFactory: instantiateFactory,
		configurations:     configurations,
	}

	f.configurations = configurations
	_ = f.SetInstance("configurations", configurations)

	f.builder = f.Builder()

	f.Append(f)
	return f
}

// SystemConfiguration getter
func (f *configurableFactory) SystemConfiguration() *system.Configuration {
	return f.systemConfig
}

// Configuration getter
func (f *configurableFactory) Configuration(name string) interface{} {
	cfg, ok := f.configurations.Get(name)
	if ok {
		return cfg
	}
	return nil
}

// BuildProperties build all properties
func (f *configurableFactory) BuildProperties() (systemConfig *system.Configuration, err error) {
	// manually inject systemConfiguration
	systemConfig = f.GetInstance(system.Configuration{}).(*system.Configuration)
	_ = f.InjectDefaultValue(systemConfig)

	profile := os.Getenv(EnvAppProfilesActive)
	if profile == "" {
		profile = defaultProfileName
	}
	f.builder.SetDefaultProperty(PropAppProfilesActive, profile)

	for prop, val := range f.DefaultProperties() {
		f.builder.SetDefaultProperty(prop, val)
	}

	_, err = f.builder.Build(profile)
	if err == nil {
		_ = f.InjectIntoObject(systemConfig)
		//replacer.Replace(systemConfig, systemConfig)

		f.configurations.Set(System, systemConfig)

		f.systemConfig = systemConfig
	}

	//load system properties
	allProperties := f.GetInstances(at.ConfigurationProperties{})
	log.Debug(len(allProperties))
	for _, properties := range allProperties {
		_ = f.builder.Load(properties.MetaObject)
	}
	return
}

// Build build all auto configurations
func (f *configurableFactory) Build(configs []*factory.MetaData) {
	// categorize configurations first, then inject object if necessary
	for _, item := range configs {
		if annotation.Contains(item.MetaObject, at.AutoConfiguration{}) {
			f.configureContainer = append(f.configureContainer, item)
		} else {
			err := ErrInvalidObjectType
			log.Errorf("item: %v err: %v", item, err)
		}
	}

	f.build(f.configureContainer)

	// load properties again
	//allProperties := f.GetInstances(at.ConfigurationProperties{})
	//log.Debug(len(allProperties))
	//for _, properties := range allProperties {
	//	_ = f.builder.Load(properties.MetaObject)
	//}
}

// Instantiate run instantiation by method
func (f *configurableFactory) Instantiate(configuration interface{}) (err error) {
	cv := reflect.ValueOf(configuration)
	icv := reflector.Indirect(cv)

	//if !cv.IsValid() {
	//	return ErrInvalidObjectType
	//}
	configType := cv.Type()
	//log.Debug("type: ", configType)
	//name := configType.Elem().Name()
	//log.Debug("fieldName: ", name)
	pkgName := io.DirName(icv.Type().PkgPath())
	var runtimeDeps factory.Deps
	rd := icv.FieldByName("RuntimeDeps")
	if rd.IsValid() {
		runtimeDeps = rd.Interface().(factory.Deps)
	}
	// call Init
	numOfMethod := cv.NumMethod()
	//log.Debug("methods: ", numOfMethod)
	for mi := 0; mi < numOfMethod; mi++ {
		// get method
		// find the dependencies of the method
		method := configType.Method(mi)
		methodName := str.LowerFirst(method.Name)
		if rd.IsValid() {
			// append inst to f.components
			deps := runtimeDeps.Get(method.Name)

			metaData := &factory.MetaData{
				Name:       pkgName + "." + methodName,
				MetaObject: method,
				DepNames:   deps,
			}
			f.AppendComponent(configuration, metaData)
		} else {
			f.AppendComponent(configuration, method)
		}
	}
	return
}

func (f *configurableFactory) parseName(item *factory.MetaData) string {

	//return item.PkgName
	name := strings.Replace(item.TypeName, PostfixConfiguration, "", -1)
	name = str.ToLowerCamel(name)

	if name == "" || name == strings.ToLower(PostfixConfiguration) {
		name = item.PkgName
	}
	return name
}

func (f *configurableFactory) injectProperties(cf interface{}) {
	v := reflect.ValueOf(cf)
	cfv := reflector.Indirect(v)
	cft := cfv.Type()
	for _, field := range reflector.DeepFields(cft) {
		var fieldObjValue reflect.Value

		// find properties field
		if !annotation.Contains(field.Type, at.ConfigurationProperties{}) {
			continue
		}

		if cfv.IsValid() && cfv.Kind() == reflect.Struct {
			fieldObjValue = cfv.FieldByName(field.Name)
		}

		// find it first
		injectedObject := f.GetInstance(field.Type)

		if !annotation.Contains(injectedObject, at.AutoWired{}) {
			continue
		}

		var injectedObjectValue reflect.Value
		if injectedObject == nil {
			injectedObjectValue = reflect.New(reflector.IndirectType(field.Type))
		} else {
			injectedObjectValue = reflect.ValueOf(injectedObject)
		}

		if fieldObjValue.CanSet() && injectedObjectValue.Type().AssignableTo(fieldObjValue.Type()) {
			fieldObjValue.Set(injectedObjectValue)
		} else {
			log.Errorf("Error: trying to assign %v to %v", injectedObjectValue.Type(), fieldObjValue.Type())
		}
	}
	return
}

func (f *configurableFactory) build(cfgContainer []*factory.MetaData) {
	var err error
	for _, item := range cfgContainer {
		name := f.parseName(item)
		config := item.MetaObject

		isContextAware := annotation.Contains(item.MetaObject, at.ContextAware{})
		if f.systemConfig != nil {
			if !isContextAware &&
				f.systemConfig != nil && !str.InSlice(name, f.systemConfig.App.Profiles.Include) {
				log.Warnf("Auto configuration %v is filtered out! Just ignore this warning if you intended to do so.", name)
				continue
			}
		}
		log.Infof("Auto configuration %v is configured on %v.", item.PkgName, item.Type)

		err = f.initProperties(config)

		// inject into func
		var cf interface{}
		if item.Kind == types.Func {
			cf, err = f.InjectIntoFunc(config)
		}
		if err == nil && cf != nil {
			// new properties
			// we have two choices: the first is current implementation which inject properties by default,
			// the second is inject properties and load to the container, let user to decide inject to configuration through constructor
			f.injectProperties(cf)

			// inject other fields
			_ = f.InjectIntoObject(cf)

			// instantiation
			_ = f.Instantiate(cf)

			// save configuration
			configName := name
			//if _, ok := f.configurations.Get(name); ok {
			//	configName = reflector.GetFullName(cf)
			//}
			// TODO: should set full name instead
			f.configurations.Set(configName, cf)
		} else {
			log.Warn(err)
		}
	}
}

func (f *configurableFactory) initProperties(config interface{}) (err error) {
	cft, ok := reflector.GetObjectType(config)
	if ok {
		// load properties
		for _, field := range reflector.DeepFields(cft) {

			// find properties field
			if !annotation.Contains(field.Type, at.ConfigurationProperties{}) {
				continue
			}
			newPropVal := reflect.New(reflector.IndirectType(field.Type))
			newPropObj := newPropVal.Interface()
			err = f.InjectDefaultValue(newPropObj)
			if err != nil {
				log.Warn(err)
				return
			}

			// load properties, inject settings
			err = f.builder.Load(newPropObj)
			if err != nil {
				log.Warn(err)
				return
			}

			// save new properties to container
			err = f.SetInstance(newPropObj)
			if err != nil {
				log.Warn(err)
				return
			}
		}
	}
	return
}

func (f *configurableFactory) StartSchedulers(schedulerServices []*factory.MetaData) (schedulers []*scheduler.Scheduler) {
	for _, svcMD := range schedulerServices {
		svc := svcMD.Instance
		methods, annotations := annotation.FindAnnotatedMethods(svc, at.Scheduled{})
		for i, method := range methods {
			sch := scheduler.NewScheduler()
			ann := annotations[i]
			_ = annotation.Inject(ann)
			schAnn := ann.Field.Value.Interface().(at.Scheduled)
			if schAnn.AtCron != nil {
				sch.RunWithExpr(schAnn.AtTag, schAnn.AtCron,
					func() {
						f.runTask(svc, method, ann, sch)
					},
				)
			} else {
				sch.Run(schAnn.AtTag, schAnn.AtLimit, schAnn.AtEvery, schAnn.AtUnit, schAnn.AtTime, schAnn.AtDelay,
					func() {
						f.runTask(svc, method, ann, sch)
					},
				)
			}

			schedulers = append(schedulers, sch)
		}
	}
	return nil
}

func (f *configurableFactory) runTask(svc interface{}, method reflect.Method, ann *annotation.Annotation, sch *scheduler.Scheduler) {
	result, err := reflector.CallMethodByName(svc, method.Name, ann.Parent.Value.Interface())
	if err == nil {
		switch result.(type) {
		case bool:
			res := result.(bool)
			if res {
				sch.Stop()
			}
		}
	}
}

