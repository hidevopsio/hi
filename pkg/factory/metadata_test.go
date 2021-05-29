package factory

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"reflect"
	"testing"
)

type foo struct {
	at.Qualifier `value:"foo"`
	name string
}

type fooBarService struct {
	HelloWorld HelloWorld `inject:""`
	foo        *foo
}

type Hello string
type HelloWorld string
type HelloHiboot string

type helloConfiguration struct {
	Configuration
}

func (c *helloConfiguration) Hello() Hello {
	return Hello("Hello World")
}

func (c *helloConfiguration) HelloWorld(h Hello) HelloWorld {
	return HelloWorld(h + "World")
}

func (c *helloConfiguration) HelloHiboot(h Hello) HelloHiboot {
	return HelloHiboot(h + "Hello Hiboot")
}

func newFooBarService(foo *foo) *fooBarService {
	return &fooBarService{foo: foo}
}

func TestUtils(t *testing.T) {
	helloConfig := new(helloConfiguration)
	helloTyp := reflect.TypeOf(helloConfig)
	numOfMethod := helloTyp.NumMethod()
	//log.Debug("methods: ", numOfMethod)
	methodTestData := make([]*MetaData, 0)
	for mi := 0; mi < numOfMethod; mi++ {
		method := helloTyp.Method(mi)
		// append inst to f.components
		methodTestData = append(methodTestData, NewMetaData(helloConfig, method))
	}

	t.Run("should parse instance name via object", func(t *testing.T) {
		md := NewMetaData("", new(fooBarService))
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, "factory", md.PkgName)
		assert.NotEqual(t, nil, md.MetaObject)
	})

	t.Run("should parse instance name via object", func(t *testing.T) {
		md := NewMetaData("", new(fooBarService))
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, "factory", md.PkgName)
		assert.NotEqual(t, nil, md.MetaObject)
	})

	t.Run("should parse instance name via object with eliminator", func(t *testing.T) {
		md := NewMetaData(new(fooBarService))
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, "factory.fooBarService", md.Name)
		assert.NotEqual(t, nil, md.MetaObject)
	})

	t.Run("should parse object instance name via constructor", func(t *testing.T) {
		md := NewMetaData("", newFooBarService)
		assert.Equal(t, "fooBarService", md.TypeName)
		assert.Equal(t, types.Func, md.Kind)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		md := NewMetaData(svc)
		assert.Equal(t, "factory.service", md.Name)
		assert.Equal(t, svc, md.MetaObject)
	})

	t.Run("should parse object instance name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		md := NewMetaData("foo", svc)
		assert.Equal(t, "service", md.TypeName)
		assert.Equal(t, "factory.foo", md.Name)
		assert.Equal(t, svc, md.MetaObject)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		md := NewMetaData(&MetaData{MetaObject: new(service)})
		assert.Equal(t, "factory.service", md.Name)
		assert.Equal(t, svc, md.MetaObject)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		name, object := ParseParams(svc)
		assert.Equal(t, "factory.service", name)
		assert.Equal(t, svc, object)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		type service struct{}
		svc := new(service)
		name, object := ParseParams("factory.fooService", svc)
		assert.Equal(t, "factory.fooService", name)
		assert.Equal(t, svc, object)
	})

	t.Run("should parse object pkg name", func(t *testing.T) {
		name, _ := ParseParams("factory.fooService")
		assert.Equal(t, "factory.fooService", name)
	})

	t.Run("should parse func dependencies", func(t *testing.T) {
		fn := newFooBarService
		ft, ok := reflector.GetObjectType(fn)
		assert.Equal(t, true, ok)

		deps := parseDependencies(fn, types.Func, ft)
		assert.Equal(t, []string{"factory.foo"}, deps)
	})

	t.Run("should append dep", func(t *testing.T) {
		dep := appendDep("", "a.b")
		dep = appendDep(dep, "c.d")
		assert.Equal(t, "a.b,c.d", dep)
	})

	t.Run("should clone meta data", func(t *testing.T) {
		src := NewMetaData(new(foo))
		dst := CloneMetaData(src)
		assert.Equal(t, dst.Type, src.Type)
	})
	t.Run("test GetObjectQualifierName", func(t *testing.T) {
		name := GetObjectQualifierName(reflect.ValueOf(new(foo)), "default")
		assert.Equal(t, "foo", name)
	})
}

func TestParseParams(t *testing.T) {
	type service1 struct{}
	type service2 struct{}
	type Service3 struct {
		service2 service2
	}

	type service4 struct {
		Service3 `depends:"factory.service2"`
	}

	svc1 := new(service1)
	svc2 := new(service2)
	svc3 := new(Service3)
	md := NewMetaData(new(service4))

	iTyp := reflector.IndirectType(reflect.TypeOf(svc3))

	testData := []struct {
		p1   interface{}
		p2   interface{}
		name string
		obj  interface{}
	}{
		{md, nil, "factory.service4", md},
		{svc1, nil, "factory.service1", svc1},
		{service1{}, svc1, "factory.service1", svc1},
		{"factory.myService", svc2, "factory.myService", svc2},
		{"myService", svc2, "factory.myService", svc2},
		{"factory.service", nil, "factory.service", nil},
		{iTyp, MetaData{}, "factory.service3", MetaData{}},
	}

	var name string
	var obj interface{}
	for _, d := range testData {
		if d.p2 == nil {
			name, obj = ParseParams(d.p1)
		} else {
			name, obj = ParseParams(d.p1, d.p2)
		}
		assert.Equal(t, name, d.name)
		assert.Equal(t, obj, d.obj)
	}
}
