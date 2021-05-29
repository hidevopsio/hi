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

package instantiate_test

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"os"
	"reflect"
	"testing"
)

const (
	helloWorld = "Hello world"
	hiboot     = "Hiboot"
)

type ContextAwareFooBar struct {
	at.ContextAware
	Name    string
	context context.Context
}

func newContextAwareFooBar(context context.Context) *ContextAwareFooBar {
	return &ContextAwareFooBar{context: context}
}

type FooBar struct {
	Name string
}

type fooBarService struct {
	Name   string
	fooBar *FooBar
}

type qualifierService struct {
	at.Qualifier `value:"instantiate_test.hibootService"`
}

type testService struct {
	at.Qualifier `value:"instantiate_test.hibootService"`
}

type BarService interface {
	Baz() string
}

type BarServiceImpl struct {
	BarService
}

func (s *BarServiceImpl) Baz() string {
	return "baz"
}

func newFooBarService(fooBar *FooBar) *fooBarService {
	return &fooBarService{
		fooBar: fooBar,
	}
}

type helloService struct {
}

func newHelloService() *helloService {
	return &helloService{}
}

func newHelloNilService() *helloService {
	return nil
}

type HelloWorld struct {
	Message string
}

func (s *helloService) HelloWorld() *HelloWorld {
	return &HelloWorld{Message: helloWorld}
}

func TestInstantiateFactory(t *testing.T) {
	type foo struct{ Name string }
	f := new(foo)

	instFactory := instantiate.NewInstantiateFactory(nil, nil, nil)
	testName := "foobar"
	t.Run("should failed to set/get instance when factory is not initialized", func(t *testing.T) {
		inst := instFactory.GetInstance("not-exist-instance")
		assert.Equal(t, nil, inst)

		err := instFactory.SetInstance("foo", nil)
		assert.Equal(t, instantiate.ErrNotInitialized, err)

		item := instFactory.Items()
		// should have 1 instance (of system.Configuration)
		assert.Equal(t, 6, len(item))
	})

	hello := newHelloService()
	helloTyp := reflect.TypeOf(hello)
	numOfMethod := helloTyp.NumMethod()
	//log.Debug("methods: ", numOfMethod)
	testComponents := make([]*factory.MetaData, 0)
	for mi := 0; mi < numOfMethod; mi++ {
		method := helloTyp.Method(mi)
		// append inst to f.components
		testComponents = append(testComponents, factory.NewMetaData(hello, method))
	}

	testComponents = append(testComponents, factory.NewMetaData(f),
		factory.NewMetaData(web.NewContext(nil)),
		factory.NewMetaData(newContextAwareFooBar),
		factory.NewMetaData(&FooBar{Name: testName}),
		factory.NewMetaData(&BarServiceImpl{}),
		factory.NewMetaData(newFooBarService),
		factory.NewMetaData(new(qualifierService)),
		factory.NewMetaData(newHelloNilService),
	)

	ic := cmap.New()
	customProps := cmap.New()
	customProps.Set("app.project", "instantiate-test")
	t.Run("should create new instantiate factory", func(t *testing.T) {
		os.Setenv(autoconfigure.EnvAppProfilesActive, "")
		f := instantiate.NewInstantiateFactory(ic, testComponents, customProps)
		assert.Equal(t, nil, f.GetProperty("app.profiles.active"))

		f.SetProperty("foo.bar", "foo bar")
		assert.Equal(t, "foo bar", f.GetProperty("foo.bar"))
	})

	customProps.Set("app.profiles.active", "local")
	instFactory = instantiate.NewInstantiateFactory(ic, testComponents, customProps)
	instFactory.AppendComponent(new(testService))
	instFactory.AppendComponent("context.context", web.NewContext(nil))
	t.Run("should initialize factory", func(t *testing.T) {
		assert.Equal(t, true, instFactory.Initialized())
	})

	t.Run("should initialize factory", func(t *testing.T) {
		cstProp := instFactory.DefaultProperties()
		assert.NotEqual(t, 0, len(cstProp))
	})

	t.Run("should build components", func(t *testing.T) {
		instFactory.BuildComponents()
	})

	t.Run("should get built instance", func(t *testing.T) {
		inst := instFactory.GetInstance(HelloWorld{})
		assert.NotEqual(t, nil, inst)
		assert.Equal(t, "Hello world", inst.(*HelloWorld).Message)
	})

	t.Run("should get built instance in specific type", func(t *testing.T) {
		hmd := instFactory.GetInstance(HelloWorld{}, factory.MetaData{})
		assert.NotEqual(t, nil, hmd)
		inst := hmd.(*factory.MetaData).Instance
		assert.Equal(t, "Hello world", inst.(*HelloWorld).Message)
	})

	t.Run("should set and get instance from factory", func(t *testing.T) {
		instFactory.SetInstance(f)
		inst := instFactory.GetInstance(foo{})
		assert.Equal(t, f, inst)
	})

	t.Run("should failed to get instance that does not exist", func(t *testing.T) {
		inst := instFactory.GetInstance("not-exist-instance")
		assert.Equal(t, nil, inst)
	})

	t.Run("should failed to get instances that does not exist", func(t *testing.T) {
		inst := instFactory.GetInstances("not-exist-instances")
		assert.Equal(t, 0, len(inst))
	})

	t.Run("should set instance", func(t *testing.T) {
		err := instFactory.SetInstance(new(foo))
		assert.NotEqual(t, nil, err)
	})

	t.Run("should get factory items", func(t *testing.T) {
		items := instFactory.Items()
		assert.NotEqual(t, 0, len(items))
	})

	t.Run("should get qualifierService with qualifier name hibootService", func(t *testing.T) {
		svc := instFactory.GetInstance("instantiate_test.hibootService")
		assert.NotEqual(t, 0, svc)
	})

	t.Run("should get appended testService", func(t *testing.T) {
		svc := instFactory.GetInstance(testService{})
		assert.NotEqual(t, 0, svc)
	})

	t.Run("should inject dependency by method InjectDependency", func(t *testing.T) {
		instFactory.InjectDependency(factory.NewMetaData(newFooBarService))
	})

	builder := instFactory.Builder()
	builder.Build()
	t.Run("should replace property", func(t *testing.T) {
		builder.SetProperty("app.name", "foo")
		builder.SetProperty("app.profiles.include", []string{"foo", "bar"})

		res := instFactory.Replace("this is ${app.name} property")
		assert.Equal(t, "this is foo property", res)

		res = instFactory.Replace("${app.profiles.include}")
		assert.Equal(t, []string{"foo", "bar"}, res)
	})

	t.Run("should replace with default property", func(t *testing.T) {
		res := instFactory.Replace("this is ${default.value:default} property")
		assert.Equal(t, "this is default property", res)
	})

	t.Run("should replace with environment variable", func(t *testing.T) {
		res := instFactory.Replace("this is ${HOME}")
		home := os.Getenv("HOME")
		assert.Equal(t, "this is "+home, res)
	})

	t.Run("should get property", func(t *testing.T) {
		res := instFactory.GetProperty("app.name")
		assert.Equal(t, "foo", res)
	})

	type Greeter struct {
		at.AutoWired

		Name string `default:"Hiboot"`
	}
	greeter := &Greeter{
		Name: "Hiboot",
	}
	t.Run("should inject default value", func(t *testing.T) {
		err := instFactory.InjectDefaultValue(greeter)
		assert.Equal(t, nil, err)
		assert.Equal(t, hiboot, greeter.Name)
	})

	type DevTester struct {
		Greeter *Greeter
		Home    string   `value:"${HOME}"`
	}
	devTester := new(DevTester)
	instFactory.SetInstance(greeter)
	t.Run("should inject into object", func(t *testing.T) {
		err := instFactory.InjectIntoObject(devTester)
		assert.Equal(t, nil, err)
		assert.Equal(t, hiboot, devTester.Greeter.Name)
		assert.Equal(t, os.Getenv("HOME"), devTester.Home)
	})

	devTesterConstructor := func(g *Greeter) *DevTester {
		return &DevTester{
			Greeter: g,
		}
	}

	t.Run("should inject into func", func(t *testing.T) {
		obj, err := instFactory.InjectIntoFunc(devTesterConstructor)
		assert.Equal(t, nil, err)
		assert.Equal(t, hiboot, obj.(*DevTester).Greeter.Name)
	})

	t.Run("should inject into method", func(t *testing.T) {
		obj, err := instFactory.InjectIntoMethod(nil, nil)
		assert.Equal(t, inject.ErrInvalidMethod, err)
		assert.Equal(t, nil, obj)
	})

	svc := newHelloService()
	t.Run("should inject into method", func(t *testing.T) {
		typ := reflect.TypeOf(svc)
		method, ok := typ.MethodByName("HelloWorld")
		assert.Equal(t, true, ok)
		obj, err := instFactory.InjectIntoMethod(svc, method)
		assert.Equal(t, nil, err)
		assert.Equal(t, helloWorld, obj.(*HelloWorld).Message)
	})
}

func TestMapSet(t *testing.T) {
	requiredClasses := mapset.NewSet()
	requiredClasses.Add("Cooking")
	requiredClasses.Add("English")
	requiredClasses.Add("Math")
	requiredClasses.Add("Biology")

	scienceSlice := []interface{}{"Biology", "Chemistry"}
	scienceClasses := mapset.NewSetFromSlice(scienceSlice)

	electiveClasses := mapset.NewSet()
	electiveClasses.Add("Welding")
	electiveClasses.Add("Music")
	electiveClasses.Add("Automotive")

	bonusClasses := mapset.NewSet()
	bonusClasses.Add("Go Programming")
	bonusClasses.Add("Python Programming")

	//Show me all the available classes I can take
	allClasses := requiredClasses.Union(scienceClasses).Union(electiveClasses).Union(bonusClasses)
	fmt.Println(allClasses)

	//Is cooking considered a science class?
	fmt.Println(scienceClasses.Contains("Cooking")) //false

	//Show me all classes that are not science classes, since I hate science.
	fmt.Println(allClasses.Difference(scienceClasses)) //Set{Music, Automotive, Go Programming, Python Programming, Cooking, English, Math, Welding}

	//Which science classes are also required classes?
	fmt.Println(scienceClasses.Intersect(requiredClasses)) //Set{Biology}

	//How many bonus classes do you offer?
	fmt.Println(bonusClasses.Cardinality()) //2

	//Do you have the following classes? Welding, Automotive and English?
	fmt.Println(allClasses.IsSuperset(mapset.NewSetFromSlice([]interface{}{"Welding", "Automotive", "English"}))) //true
}

type contextAwareFuncObject struct {
	at.ContextAware

	context context.Context
}
type contextAwareMethodObject struct {
	at.ContextAware

	context context.Context
}

func newContextAwareObject(ctx context.Context) *contextAwareFuncObject {
	//log.Infof("context: %v", ctx)
	return &contextAwareFuncObject{context: ctx}
}

type foo struct {

}

func (f *foo) ContextAwareMethodObject(ctx context.Context) *contextAwareMethodObject {
	return &contextAwareMethodObject{context: ctx}
}

func TestRuntimeInstance(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	//log.Debug("methods: ", numOfMethod)
	testComponents := make([]*factory.MetaData, 0)

	ctx := web.NewContext(nil)
	// method
	f := new(foo)
	ft := reflect.TypeOf(f)

	ctxMd := factory.NewMetaData(reflector.GetLowerCamelFullName(new(context.Context)), ctx)
	method, ok := ft.MethodByName("ContextAwareMethodObject")
	assert.Equal(t, true, ok)
	testComponents = append(testComponents,
		factory.NewMetaData(f, method),
		ctxMd,
		factory.NewMetaData(newContextAwareObject),
	)

	ic := cmap.New()
	customProps := cmap.New()
	customProps.Set("app.project", "runtime-test")
	instFactory := instantiate.NewInstantiateFactory(ic, testComponents, customProps)
	instFactory.AppendComponent(new(testService))
	_ = instFactory.BuildComponents()
	dps := instFactory.GetInstances(at.ContextAware{})
	if len(dps) > 0 {
		ri, err := instFactory.InjectContextAwareObjects(web.NewContext(nil), dps)
		assert.Equal(t, nil, err)
		log.Debug(ri.Items())
		assert.Equal(t, ctx, ri.Get(new(context.Context)))
		assert.NotEqual(t, nil, ri.Get(contextAwareFuncObject{}))
		assert.NotEqual(t, nil, ri.Get(contextAwareMethodObject{}))
	}
}
