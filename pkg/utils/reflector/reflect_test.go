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

package reflector

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector/tester"
	"reflect"
	"testing"
)

type Foo struct {
	Name     string
	Age      int
	nickname string
}

func newFoo(name string) *Foo {
	return &Foo{Name: name}
}

func (f *Foo) Bar() {

}

func (f *Foo) Nickname() string {
	return f.nickname
}

func (f *Foo) SetNickname(nickname string) *Foo {
	f.nickname = nickname
	return f
}

type Bar struct {
	Name string
	Age  int
}

type Path string

type Baz struct {
	Path `value:"test"`
	Foo
	Bar Bar
}

type FooBar struct{}

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestLookupField(t *testing.T) {
	t.Run("should find anonymous field on Baz", func(t *testing.T) {
		assert.Equal(t, true, HasField(&Baz{}, "Foo"))
	})

	t.Run("should not find anonymous field on Foo", func(t *testing.T) {
		assert.Equal(t, false, HasField(&Bar{}, "Foo"))
	})
}

//func TestConstructor(t *testing.T) {
//	objVal := reflect.ValueOf(&Foo{})
//
//	objType := objVal.Type()
//	log.Debug("type: ", objType)
//	objName := objType.Elem().Name()
//	log.Debug("name: ", objName)
//
//	object := objVal.Interface()
//	log.Debug("object: ", object)
//
//	// call Init
//	method, ok := objType.MethodByName("Init")
//	if ok {
//		methodType := method.Type
//		numIn := methodType.NumIn()
//		inputs := make([]reflect.Value, numIn)
//		for i := 0; i < numIn; i++ {
//			t := methodType.In(i)
//
//			log.Debugf("%v: %v %v", i, t.Name(), t)
//		}
//		inputs[0] = reflect.ValueOf(object)
//		//method.Func.Call(inputs)
//	}
//}

func TestNewReflectType(t *testing.T) {
	foo := NewReflectType(Foo{})
	assert.NotEqual(t, nil, foo)
}

func TestValidate(t *testing.T) {
	foo := &Foo{Name: "foo"}

	t.Run("should validate object foo", func(t *testing.T) {
		f, err := Validate(foo)
		assert.Equal(t, nil, err)
		assert.Equal(t, reflect.Struct, f.Kind())
	})

	t.Run("should return error that value is unaddressable", func(t *testing.T) {
		_, err := Validate(123)
		assert.Equal(t, "value is unaddressable", err.Error())
	})

	t.Run("should return error that value is not valid", func(t *testing.T) {
		_, err := Validate((*Foo)(nil))
		assert.Equal(t, "value is not valid", err.Error())
	})
}

func TestDeepFields(t *testing.T) {
	baz := &Baz{Bar: Bar{Name: "bar"}}
	baz.Name = "foo"
	bt := reflect.TypeOf(baz)
	df := DeepFields(bt)
	assert.Equal(t, 5, len(df))
	assert.Equal(t, "Path", df[0].Name)
	assert.Equal(t, "Name", df[1].Name)
	assert.Equal(t, "Age", df[2].Name)
}

func TestIndirect(t *testing.T) {
	foo := &Foo{Name: "foo"}
	f := reflect.ValueOf(foo)
	fv := Indirect(f)
	assert.Equal(t, reflect.Ptr, f.Kind())
	assert.Equal(t, reflect.Struct, fv.Kind())
}

func TestIndirectValue(t *testing.T) {
	foo := Foo{Name: "foo"}
	fv := IndirectValue(foo)
	assert.Equal(t, reflect.Struct, fv.Kind())

	fv = IndirectValue(&foo)
	assert.Equal(t, reflect.Struct, fv.Kind())

	fv = IndirectValue("test")
	assert.Equal(t, reflect.String, fv.Kind())

	fv = IndirectValue(reflect.ValueOf(foo))
	assert.Equal(t, reflect.Struct, fv.Kind())

	fv = IndirectValue(reflect.TypeOf(foo))
	assert.Equal(t, false, fv.IsValid())
}

func TestIndirectType(t *testing.T) {
	foo := &Foo{Name: "foo"}
	f := reflect.TypeOf(foo)
	ft := IndirectType(f)
	assert.Equal(t, reflect.Ptr, f.Kind())
	assert.Equal(t, reflect.Struct, ft.Kind())
}

func TestGetFieldValue(t *testing.T) {
	foo := &Foo{Name: "foo"}
	fv := GetFieldValue(foo, "Name")
	assert.Equal(t, "foo", fv.Interface())
}

func TestFindFieldByTag(t *testing.T) {
	type foo struct {
		Name string `json:"name"`
		AtValues []string `at:"values"`
	}
	f := &foo{}
	field, ok := FindFieldByTag(f, "at", "values")
	assert.Equal(t, true, ok)
	assert.Equal(t, "AtValues", field.Name)
}

func TestSetFieldValue(t *testing.T) {
	foo := &Foo{}
	value := "foo"
	t.Run("should set field value", func(t *testing.T) {
		// set field object
		err := SetFieldValue(foo, "Name", value)
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", value)
	})

	t.Run("should not set invalid object", func(t *testing.T) {
		x := 123
		err := SetFieldValue(x, "Name", 321)
		assert.Equal(t, ErrInvalidInput, err)
	})

	t.Run("should not set invalid object", func(t *testing.T) {
		err := SetFieldValue((*Foo)(nil), "Name", value)
		assert.Equal(t, ErrInvalidInput, err)
	})

	t.Run("should not set invalid object", func(t *testing.T) {
		err := SetFieldValue(foo, "nickname", value)
		assert.Equal(t, ErrFieldCanNotBeSet, err)
	})
}

func TestGetKind(t *testing.T) {

	t.Run("should return reflect.Unit for uint64", func(t *testing.T) {
		var x uint64
		x = 1234
		k := GetKindByValue(reflect.ValueOf(x))
		assert.Equal(t, reflect.Uint, k)
	})

	t.Run("should return reflect.Int for int64", func(t *testing.T) {
		var x int64
		x = 1234
		k := GetKindByValue(reflect.ValueOf(x))
		assert.Equal(t, reflect.Int, k)
	})

	t.Run("should return reflect.Float32 for float64", func(t *testing.T) {
		var x float64
		x = 1.234
		k := GetKindByValue(reflect.ValueOf(x))
		assert.Equal(t, reflect.Float32, k)
	})

	t.Run("should return Ptr", func(t *testing.T) {
		k := GetKindByValue(reflect.ValueOf((*Foo)(nil)))
		assert.Equal(t, reflect.Ptr, k)
	})

	t.Run("should return Int", func(t *testing.T) {
		k := GetKindByType(reflect.ValueOf(int(1)).Type())
		assert.Equal(t, reflect.Int, k)
	})

	t.Run("should return Uint", func(t *testing.T) {
		k := GetKindByType(reflect.ValueOf(uint(1)).Type())
		assert.Equal(t, reflect.Uint, k)
	})

	t.Run("should return Bool", func(t *testing.T) {
		k := GetKindByType(reflect.ValueOf(true).Type())
		assert.Equal(t, reflect.Bool, k)
	})

	t.Run("should return Float32", func(t *testing.T) {
		k := GetKindByType(reflect.ValueOf(0.01).Type())
		assert.Equal(t, reflect.Float32, k)
	})

	t.Run("should return String", func(t *testing.T) {
		k := GetKindByType(reflect.ValueOf("abc").Type())
		assert.Equal(t, reflect.String, k)
	})
}

func TestValidateReflectType(t *testing.T) {
	baz := &Baz{Foo: Foo{Name: "foo"}, Bar: Bar{Name: "bar"}}

	t.Run("should validete reflect type", func(t *testing.T) {
		ValidateReflectType(baz, func(value *reflect.Value, reflectType reflect.Type, fieldSize int, isSlice bool) error {
			assert.Equal(t, *baz, value.Interface())
			assert.Equal(t, 1, fieldSize)
			assert.Equal(t, "Baz", reflectType.Name())
			assert.Equal(t, false, isSlice)
			return nil
		})
	})

	t.Run("should return value is not valid", func(t *testing.T) {
		err := ValidateReflectType((*Foo)(nil), nil)
		assert.Equal(t, "value is not valid", err.Error())
	})

	t.Run("should validate slice", func(t *testing.T) {
		s := []int{1, 2, 3}
		err := ValidateReflectType(&s, nil)
		assert.Equal(t, nil, err)
	})

}

func TestGetName(t *testing.T) {

	t.Run("should get struct name by pointer", func(t *testing.T) {
		n := GetName(new(Foo))
		assert.Equal(t, "Foo", n)
	})

	t.Run("should get struct name", func(t *testing.T) {
		n := GetName(Foo{})
		assert.Equal(t, "Foo", n)
	})

	t.Run("should get struct name by lower case", func(t *testing.T) {
		n := GetLowerCamelName(new(Foo))
		assert.Equal(t, "foo", n)
	})

	t.Run("should return InvalidInputError if input is nil", func(t *testing.T) {
		n := GetLowerCamelName((*Foo)(nil))
		assert.Equal(t, "", n)
	})
}

func TestCallMethodByName(t *testing.T) {
	foo := new(Foo)
	t.Run("should call method SetNickname", func(t *testing.T) {
		res, err := CallMethodByName(foo, "SetNickname", "foobar")
		assert.Equal(t, nil, err)
		assert.Equal(t, foo, res)
		assert.Equal(t, "foobar", foo.nickname)
	})

	t.Run("should call method Nickname", func(t *testing.T) {
		res, err := CallMethodByName(foo, "Nickname")
		assert.Equal(t, nil, err)
		assert.Equal(t, "foobar", res)
	})

	t.Run("should call method Bar", func(t *testing.T) {
		res, err := CallMethodByName(foo, "Bar")
		assert.Equal(t, nil, err)
		assert.Equal(t, nil, res)
	})

	t.Run("should call method by name", func(t *testing.T) {
		_, err := CallMethodByName(foo, "NotExistMethod")
		assert.Equal(t, ErrInvalidMethod, err)
	})

}

type Embedded struct{}
type PtrEmbedded struct{}

type hidden struct{}

type Test struct {
	Embedded
	hidden
	*PtrEmbedded
	Exportedvalue int
	unexported    string
	A             struct{ A int }
}

func TestEmbedded(t *testing.T) {
	type fooInterface interface {
	}

	testCases := []struct {
		name      string
		source    interface{}
		fieldName string
		ok        bool
		private   bool
		anonymous bool
	}{
		{
			name:      "Should find Embedded field by name",
			source:    Test{},
			fieldName: "Embedded",
			ok:        true,
			private:   false,
			anonymous: true,
		},
		{
			name:      "Should find Embedded field by name",
			source:    &Test{},
			fieldName: "Embedded",
			ok:        true,
			private:   false,
			anonymous: true,
		},
		{
			name:      "Should find hidden field by name",
			source:    Test{},
			fieldName: "hidden",
			ok:        true,
			private:   true,
			anonymous: true,
		},
		{
			name:      "Should find PtrEmbedded field by name",
			source:    Test{},
			fieldName: "PtrEmbedded",
			ok:        true,
			private:   false,
			anonymous: true,
		},
		{
			name:      "Should find Exportedvalue field by name",
			source:    Test{},
			fieldName: "Exportedvalue",
			ok:        true,
			private:   false,
			anonymous: false,
		},
		{
			name:      "Should find unexported field by name",
			source:    Test{},
			fieldName: "unexported",
			ok:        true,
			private:   true,
			anonymous: false,
		},
		{
			name:      "Should find A field by name",
			source:    Test{},
			fieldName: "A",
			ok:        true,
			private:   false,
			anonymous: false,
		},
		{
			name:      "Should fail to parse non-struct source type",
			source:    123,
			fieldName: "A",
			ok:        false,
			private:   false,
			anonymous: false,
		},
	}

	for _, testCase := range testCases {
		typ := IndirectType(reflect.TypeOf(testCase.source))
		if typ.Kind() == reflect.Struct {
			t.Run(testCase.name, func(t *testing.T) {
				field, ok := typ.FieldByName(testCase.fieldName)
				exported := field.PkgPath != ""
				assert.Equal(t, testCase.ok, ok)
				assert.Equal(t, testCase.private, exported)
				assert.Equal(t, testCase.anonymous, field.Anonymous)
				assert.Equal(t, testCase.fieldName, field.Name)
			})
		}

		t.Run(testCase.name, func(t *testing.T) {
			ok := HasEmbeddedField(testCase.source, testCase.fieldName)
			assert.Equal(t, testCase.anonymous, ok)
		})
	}
}

func TestParseObjectName(t *testing.T) {
	t.Run("should parse object name", func(t *testing.T) {
		name := ParseObjectName(new(FooBar), "Bar")
		assert.Equal(t, "foo", name)
	})
}

func TestGetPkgPath(t *testing.T) {
	t.Run("should get object pkg path", func(t *testing.T) {
		pkgPath := GetPkgPath(Foo{})
		assert.Contains(t, "github.com/hidevopsio/hiboot/pkg/utils/reflector", pkgPath)
	})
}

func TestParseObjectPkgName(t *testing.T) {
	pkgName := ParseObjectPkgName(Foo{})
	assert.Equal(t, "reflector", pkgName)

	pkgName = ParseObjectPkgName(&Foo{})
	assert.Equal(t, "reflector", pkgName)
}

func SayHello(name string) string {
	return "Hello " + name
}

func Dummy() {
	// for test only, do nothing
}

func TestCallFunc(t *testing.T) {
	t.Run("should call func", func(t *testing.T) {
		res, err := CallFunc(SayHello, "Steve")
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, res)
		assert.Equal(t, "Hello Steve", res.(string))
	})

	t.Run("should call func", func(t *testing.T) {
		res, err := CallFunc(Dummy)
		assert.Equal(t, nil, err)
		assert.Equal(t, nil, res)
	})

	t.Run("should call func", func(t *testing.T) {
		_, err := CallFunc(int(1))
		assert.Equal(t, ErrInvalidFunc, err)
	})
}

type fooBarInterface interface {
}

type fooBarService struct {
	fooBarInterface
}

func newFooBarService() *fooBarService {
	return new(fooBarService)
}

func TestGetEmbeddedInterfaceField(t *testing.T) {
	type fakeInterface interface{}
	type fakeService struct{ fakeInterface }
	type fakeChildService struct{ fakeService }

	t.Run("should return error if on nil ", func(t *testing.T) {
		field := GetEmbeddedField(nil, "")
		assert.Equal(t, nil, field.Type)
	})

	t.Run("should return error if on nil ", func(t *testing.T) {
		field := GetEmbeddedField(newFooBarService, "fooBarInterface")
		assert.Equal(t, "fooBarInterface", field.Type.Name())
	})

	t.Run("should get embedded fakeInterface", func(t *testing.T) {
		field := GetEmbeddedField(new(fakeService), "fakeInterface")
		assert.Equal(t, "fakeInterface", field.Name)
	})

	t.Run("should get embedded fakeInterface", func(t *testing.T) {
		field := GetEmbeddedField(new(fakeChildService), "")
		assert.Equal(t, "fakeInterface", field.Name)
	})

	type fooService struct{}
	type foobarService struct{ fooService }
	t.Run("should get embedded fakeInterface", func(t *testing.T) {
		field := GetEmbeddedField(new(foobarService), "fooService", reflect.Struct)
		assert.Equal(t, "fooService", field.Name)
	})

	type FooServiceInterface interface {
	}
	type FooService struct {
	}

	t.Run("should not get embedded interface that it's not exist", func(t *testing.T) {
		type barService struct {
			FooService
		}
		field := GetEmbeddedField(new(barService), "FooService")
		assert.Equal(t, false, field.Anonymous)
	})

	t.Run("should get", func(t *testing.T) {
		type barService struct {
			FooServiceInterface `name:"foo"`
		}
		tag, ok := FindEmbeddedFieldTag(new(barService), "FooServiceInterface", "name")
		assert.Equal(t, ok, true)
		assert.Equal(t, "foo", tag)
	})

	t.Run("should get the object name", func(t *testing.T) {
		name := GetFullName(Foo{})
		assert.Equal(t, "reflector.Foo", name)
	})

	t.Run("should get the object name by object pointer", func(t *testing.T) {
		name := GetFullName(&Foo{})
		assert.Equal(t, "reflector.Foo", name)
	})

	t.Run("should get the object name by object function", func(t *testing.T) {
		name := GetFullName(newFooBarService)
		assert.Equal(t, "reflector.fooBarService", name)
	})

	t.Run("should get the object name by object pointer", func(t *testing.T) {
		name := GetLowerCamelFullNameByType(reflect.TypeOf(&Foo{}))
		assert.Equal(t, "reflector.foo", name)
	})

	t.Run("should get the object name by object pointer", func(t *testing.T) {
		name := GetFullName(&Foo{})
		assert.Equal(t, "reflector.Foo", name)
	})

	t.Run("should get the object name by object pointer", func(t *testing.T) {
		name := GetLowerCamelFullName(&Foo{})
		assert.Equal(t, "reflector.foo", name)
	})

	t.Run("should parse instance name via object", func(t *testing.T) {
		pkgName, name := GetPkgAndName(new(fooBarService))
		assert.Equal(t, "reflector", pkgName)
		assert.Equal(t, "fooBarService", name)
	})

	t.Run("should parse instance name via object", func(t *testing.T) {
		pkgName, name := GetPkgAndName(newFooBarService)
		assert.Equal(t, "reflector", pkgName)
		assert.Equal(t, "fooBarService", name)
	})

	t.Run("should get object type", func(t *testing.T) {
		foo := new(Foo)
		expectedTyp := reflect.TypeOf(foo)
		typ, ok := GetObjectType(expectedTyp)
		assert.Equal(t, true, ok)
		assert.Equal(t, "Foo", typ.Name())
	})

	t.Run("should get object type", func(t *testing.T) {
		foo := new(Foo)
		expectedTyp := reflect.ValueOf(foo)
		typ, ok := GetObjectType(expectedTyp)
		assert.Equal(t, true, ok)
		assert.Equal(t, "Foo", typ.Name())
	})

	t.Run("should get object type", func(t *testing.T) {
		_, ok := GetObjectType(nil)
		assert.Equal(t, false, ok)
	})

	t.Run("should get method out type", func(t *testing.T) {
		foo := new(Foo)
		expectedTyp := reflect.TypeOf(foo)
		method, ok := expectedTyp.MethodByName("SetNickname")
		assert.Equal(t, true, ok)
		typ, ok := GetObjectType(method)
		assert.Equal(t, true, ok)
		assert.Equal(t, "Foo", typ.Name())
	})

	t.Run("should report empty method out type", func(t *testing.T) {
		foo := new(Foo)
		expectedTyp := reflect.TypeOf(foo)
		method, ok := expectedTyp.MethodByName("Bar")
		assert.Equal(t, true, ok)
		typ, ok := GetObjectType(method)
		assert.Equal(t, false, ok)
		assert.Equal(t, nil, typ)
	})

	t.Run("should check valid object", func(t *testing.T) {
		assert.Equal(t, true, IsValidObjectType(&Foo{}))
	})

	t.Run("should check invalid object", func(t *testing.T) {
		assert.Equal(t, false, IsValidObjectType(1))
	})

	t.Run("should append component", func(t *testing.T) {
		assert.Equal(t, false, IsValidObjectType(1))
	})

	t.Run("should get specific embedded type", func(t *testing.T) {
		type EmbeddedInterfaceA interface {
		}

		type EmbeddedInterfaceB interface {
		}

		type EmbeddedString string

		type embeddedTypeA struct {
			EmbeddedInterfaceA
			EmbeddedString `value:"Hello"`
			EmbeddedInterfaceB

		}

		type embeddedTypeB struct {
			embeddedTypeA
		}
		yes := HasEmbeddedFieldType(new(embeddedTypeA), new(EmbeddedInterfaceA))
		assert.Equal(t, true, yes)

		yes = HasEmbeddedFieldType(new(embeddedTypeA), new(EmbeddedInterfaceB))
		assert.Equal(t, true, yes)

		yes = HasEmbeddedFieldType(new(embeddedTypeB), new(EmbeddedInterfaceA))
		assert.Equal(t, true, yes)

		yes = HasEmbeddedFieldType(new(embeddedTypeB), new(EmbeddedInterfaceB))
		assert.Equal(t, true, yes)

		f, ok := GetEmbeddedFieldType(new(embeddedTypeB), new(EmbeddedString))
		assert.Equal(t, true, ok)
		log.Debug(f)

		o := new(embeddedTypeB)
		field, ok := GetEmbeddedFieldByType(IndirectType(reflect.TypeOf(o)), embeddedTypeA{}, reflect.Struct)
		assert.Equal(t, true, ok)
		assert.Equal(t, field.Name, "embeddedTypeA")
	})

	t.Run("should return false if input nil on HasEmbeddedFieldType", func(t *testing.T) {
		yes := HasEmbeddedFieldType(newGreeter, new(EmbeddedAnnotation))
		assert.Equal(t, true, yes)
	})

	t.Run("should return false if input nil on HasEmbeddedFieldType", func(t *testing.T) {
		yes := HasEmbeddedFieldType(nil, nil)
		assert.Equal(t, false, yes)
	})

	t.Run("should get embedded types", func(t *testing.T) {
		type EmbedInterfaceA interface{}
		type EmbedInterfaceB interface{}
		type EmbedInterfaceC interface{}

		type EmbedStruct struct {
			EmbedInterfaceC
		}

		type MyType struct {
			EmbedInterfaceA
			EmbedInterfaceB
			EmbedStruct
		}

		embeddedTypes := GetEmbeddedFields(new(MyType))
		assert.Equal(t, 3, len(embeddedTypes))

		embeddedStructTypes := GetEmbeddedFields(new(MyType), reflect.Struct)
		assert.Equal(t, 1, len(embeddedStructTypes))
	})

	t.Run("should return false if input nil on HasEmbeddedFieldType", func(t *testing.T) {
		embeddedTypes := GetEmbeddedFields(nil)
		assert.Equal(t, 0, len(embeddedTypes))
	})

	t.Run("should return false if input nil on GetEmbeddedFieldsByType", func(t *testing.T) {
		embeddedTypes := GetEmbeddedFieldsByType(nil)
		assert.Equal(t, 0, len(embeddedTypes))
	})
}

type EmbeddedAnnotation interface {
}

type Greeter interface {
	Hello(name string) string
}

type greeter struct {
	EmbeddedAnnotation
}

func (g greeter) Hello(name string) string {
	return "Hello world"
}

func newGreeter() *greeter {
	return &greeter{}
}

func TestTypeSwitch(t *testing.T) {
	var typ interface{}
	type foo struct{}

	typ = &foo{}

	switch tp := typ.(type) {
	default:
		fmt.Printf("unexpected type %T\n", tp) // %T prints whatever type t has
	case bool:
		fmt.Printf("boolean %t\n", tp) // t has type bool
	case int:
		fmt.Printf("integer %d\n", tp) // t has type int
	case *bool:
		fmt.Printf("pointer to boolean %t\n", *tp) // t has type *bool
	case *int:
		fmt.Printf("pointer to integer %d\n", *tp) // t has type *int
	case foo:
		fmt.Printf("foo %v\n", tp) // t has type *int
	case *foo:
		fmt.Printf("pointerr to foo %v\n", *tp) // t has type *int
	}
}

func foo() {
}

type fooService struct{}

func (s *fooService) foobar() {

}

func TestGetFuncName(t *testing.T) {
	t.Run("should get function name", func(t *testing.T) {
		name := GetFuncName(foo)
		assert.Equal(t, "foo", name)
	})

	t.Run("should get method name", func(t *testing.T) {
		s := fooService{}
		name := GetFuncName(s.foobar)
		assert.Equal(t, "foobar", name)
	})

	t.Run("should get method name directly", func(t *testing.T) {
		name := GetFuncName((*fooService).foobar)
		assert.Equal(t, "foobar", name)
	})

	t.Run("should get method name by pointer", func(t *testing.T) {
		ps := &fooService{}
		name := GetFuncName(ps.foobar)
		assert.Equal(t, "foobar", name)
	})

	t.Run("should parse function name", func(t *testing.T) {
		name := "github.com/hidevopsio/hiboot/pkg/utils/reflector.(*fooService).(github.com/hidevopsio/hiboot/pkg/utils/reflector.foobar)-fm"
		name = parseFuncName(name)
		assert.Equal(t, "foobar", name)
	})

}

func TestHasEmbeddedFieldByInterface(t *testing.T) {
	type fooInterface interface {
	}

	type foo struct {
		fooInterface
	}

	t.Run("should find embedded field by type", func(t *testing.T) {
		ok := HasEmbeddedFieldType(new(foo), new(fooInterface))
		assert.Equal(t, true, ok)
	})

	type bar struct {
		foo
	}

	t.Run("should find nested embedded field by type", func(t *testing.T) {
		ok := HasEmbeddedFieldType(new(bar), new(fooInterface))
		assert.Equal(t, true, ok)
	})

	t.Run("should find nested embedded field by type with different package", func(t *testing.T) {
		type Foo struct {
			tester.Foo
		}

		type foobar struct {
			Foo
		}
		ok := HasEmbeddedFieldType(new(foobar), new(tester.Foo))
		assert.Equal(t, true, ok)
	})
}

type Model interface {
	m()
}

type Company struct{}

func (c *Company) m() {}

type MyNewCompany struct {
	Company

	Name string
}

func TestImplements(t *testing.T) {

	c := new(MyNewCompany)

	ok := Implements(c, new(Model))
	assert.True(t, ok)

}
