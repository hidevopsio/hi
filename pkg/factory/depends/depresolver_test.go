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

package depends_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/depends"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/bar"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/fake"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/foo"
	"github.com/hidevopsio/hiboot/pkg/log"
	"reflect"
	"testing"
)

type Hello string
type HelloWorld string
type HelloHiboot string

type helloConfiguration struct {
	app.Configuration
}

func (c *helloConfiguration) Hello() Hello {
	return Hello("Hello World")
}

func (c *helloConfiguration) HelloWorld(h Hello) HelloWorld {
	return HelloWorld(h + "World")
}

func (c *helloConfiguration) HelloHibootWorld(h Hello) HelloWorld {
	return HelloWorld(h + "Hiboot world")
}

func (c *helloConfiguration) HelloHiboot(h Hello) HelloHiboot {
	return HelloHiboot(h + "Hello Hiboot")
}

type helloService struct {
	HelloWorld  HelloWorld  `inject:""`
	HibootWorld HelloWorld  `inject:"helloHibootWorld"`
	HelloHiboot HelloHiboot `inject:""`
}

type fooConfiguration struct {
	app.Configuration
}

type barConfiguration struct {
	app.Configuration
}

type childConfiguration struct {
	app.Configuration `depends:"parentConfiguration"`
}

type parentConfiguration struct {
	app.Configuration `depends:"grantConfiguration"`
}

type grantConfiguration struct {
	app.Configuration `depends:"fake.configuration"`
}

type circularChildConfiguration struct {
	app.Configuration `depends:"circularParentConfiguration"`
}

type circularParentConfiguration struct {
	app.Configuration `depends:"circularGrantConfiguration"`
}

type circularGrantConfiguration struct {
	app.Configuration `depends:"circularParentConfiguration"`
}

type circularChildConfiguration2 struct {
	app.Configuration `depends:"circularParentConfiguration2"`
}

type circularParentConfiguration2 struct {
	app.Configuration `depends:"circularGrantConfiguration2"`
}

type circularGrantConfiguration2 struct {
	app.Configuration `depends:"circularChildConfiguration2"`
}

type Foo struct {
	Name string
}

type Bar struct {
	Name string
}

type Baz struct {
	Name string
}

type fooService struct {
	foo *Foo
}

func newFooService(foo *Foo) *fooService {
	return &fooService{foo: foo}
}

type barService struct {
	bar        *Bar
	fooService *fooService
}

func newBarService(bar *Bar, fooService *fooService) *barService {
	return &barService{bar: bar, fooService: fooService}
}

type foobarService struct {
	fooService *fooService
}

func newFoobarService(fooService *fooService) *foobarService {
	return &foobarService{fooService: fooService}
}

type bazService struct {
	bar *Bar
	baz *Baz
}

func newBazService(bar *Bar, baz *Baz) *bazService {
	return &bazService{bar: bar, baz: baz}
}

func TestSort(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	helloConfig := new(helloConfiguration)
	helloTyp := reflect.TypeOf(helloConfig)
	numOfMethod := helloTyp.NumMethod()
	//log.Debug("methods: ", numOfMethod)
	methodTestData := make([]*factory.MetaData, 0)
	for mi := 0; mi < numOfMethod; mi++ {
		method := helloTyp.Method(mi)
		// append inst to f.components
		methodTestData = append(methodTestData, factory.NewMetaData(helloConfig, method))
	}
	methodTestData = append(methodTestData, factory.NewMetaData(helloConfig, new(helloService)))

	testData := []struct {
		title          string
		configurations []*factory.MetaData
		err            error
	}{
		{
			title:          "should sort with method's dependencies",
			configurations: methodTestData,
		},
		{
			title: "should sort dependencies",
			configurations: []*factory.MetaData{
				factory.NewMetaData(new(bar.Configuration)),
				factory.NewMetaData(new(foo.Configuration)),
				factory.NewMetaData(new(fake.Configuration)),
				factory.NewMetaData(new(parentConfiguration)),
				factory.NewMetaData(new(grantConfiguration)),
				factory.NewMetaData(new(childConfiguration)),
				factory.NewMetaData(foo.NewConfiguration),
			},
			err: nil,
		},
		{
			title: "should sort dependencies",
			configurations: []*factory.MetaData{
				factory.NewMetaData(new(fake.Configuration)),
				factory.NewMetaData(new(fooConfiguration)),
				factory.NewMetaData(new(bar.Configuration)),
				factory.NewMetaData(new(childConfiguration)),
				factory.NewMetaData(new(grantConfiguration)),
				factory.NewMetaData(new(parentConfiguration)),
				factory.NewMetaData(foo.NewConfiguration),
				factory.NewMetaData(new(barConfiguration)),
			},
			err: nil,
		},
		{
			title: "should report some of the dependencies are not found",
			configurations: []*factory.MetaData{
				factory.NewMetaData(new(fooConfiguration)),
				factory.NewMetaData(new(childConfiguration)),
				factory.NewMetaData(new(grantConfiguration)),
				factory.NewMetaData(new(parentConfiguration)),
				factory.NewMetaData(new(barConfiguration)),
			},
			err: depends.ErrCircularDependency,
		},
		{
			title: "should sort with constructor's dependencies",
			configurations: []*factory.MetaData{
				factory.NewMetaData(newBarService),
				factory.NewMetaData(new(Bar)),
				factory.NewMetaData(newFooService),
				factory.NewMetaData(new(Foo)),
				factory.NewMetaData(new(Baz)),
				factory.NewMetaData(newBazService),
			},
			err: nil,
		},
		{
			title: "should fail to sort with circular dependencies 1",
			configurations: []*factory.MetaData{
				factory.NewMetaData(new(circularChildConfiguration)),
				factory.NewMetaData(new(circularParentConfiguration)),
				factory.NewMetaData(new(circularGrantConfiguration)),
			},
			err: depends.ErrCircularDependency,
		},
		{
			title: "should fail to sort with circular dependencies 2",
			configurations: []*factory.MetaData{
				factory.NewMetaData(new(Bar)),
				factory.NewMetaData(new(circularChildConfiguration)),
				factory.NewMetaData(new(circularParentConfiguration)),
				factory.NewMetaData(new(circularGrantConfiguration)),
			},
			err: depends.ErrCircularDependency,
		},
		{
			title: "should fail to sort with circular dependencies 3",
			configurations: []*factory.MetaData{
				factory.NewMetaData(new(circularChildConfiguration2)),
				factory.NewMetaData(new(circularParentConfiguration2)),
				factory.NewMetaData(new(circularGrantConfiguration2)),
			},
			err: depends.ErrCircularDependency,
		},
	}

	for _, data := range testData {
		t.Run(data.title, func(t *testing.T) {
			_, err := depends.Resolve(data.configurations)
			assert.Equal(t, err, data.err)
		})
	}
}
