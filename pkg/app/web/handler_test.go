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
	"github.com/kataras/iris"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
	"reflect"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

type fooController struct {
	Controller
}

// PutByIdNameAge PUT /foo/{id}/{name}/{age}
func (c *fooController) PutByIdNameAge(id int, name string, age int) error {
	log.Debugf("FooController.Put %v %v %v", id, name, age)
	return nil
}

type fakeFactory struct {
	factory.ConfigurableFactory
}

func (f *fakeFactory) GetInstance(params ...interface{}) (retVal interface{}) {
	return nil
}

func TestParse(t *testing.T) {
	restCtl := new(injectableObject)
	restCtl.object = new(fooController)
	ctrlVal := reflect.ValueOf(restCtl.object)
	method, ok := ctrlVal.Type().MethodByName("PutByIdNameAge")
	assert.Equal(t, true, ok)
	restMethod := &injectableMethod{
		requestMapping: &requestMapping{
			Method: "GET",
			Value:  "/foo/{id}/{name}/{age}",
		},
		method: &method,
	}
	hdl := newHandler(new(fakeFactory), restCtl, restMethod, at.HttpMethod{})

	t.Run("should parse method with path variable", func(t *testing.T) {
		log.Debug(hdl)
		assert.Equal(t, 3, len(hdl.pathVariable))
		assert.Equal(t, "fooController", hdl.requests[0].typeName)
		assert.Equal(t, "int", hdl.requests[1].typeName)
		assert.Equal(t, "string", hdl.requests[2].typeName)
		assert.Equal(t, "int", hdl.requests[3].typeName)
	})
	t.Run("should return ErrCanNotInterface for private field", func(t *testing.T) {
		type hello struct {
			world string
		}
		err := hdl.responseData(NewContext(iris.New()), 1, []reflect.Value{reflect.ValueOf(&hello{}).Elem().Field(0)})
		assert.Equal(t, ErrCanNotInterface, err)
	})

	t.Run("should clean path", func(t *testing.T) {
		p := clean("///a///b//c/d//e/////f/")
		assert.Equal(t, "/a/b/c/d/e/f", p)
	})

	t.Run("should clean path", func(t *testing.T) {
		p := clean("//abc/")
		assert.Equal(t, "/abc", p)
	})

	t.Run("should clean path", func(t *testing.T) {
		p := clean("//")
		assert.Equal(t, "/", p)
	})


}
