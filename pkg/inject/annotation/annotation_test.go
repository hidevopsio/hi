package annotation_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"reflect"
	"testing"
)

type AtBaz struct {
	at.Annotation

	at.BaseAnnotation
	AtCode int `value:"200" at:"code" json:"-"`
}

type AtFoo struct {
	at.Annotation

	at.BaseAnnotation
	AtID int `at:"fooId" json:"-"`
	AtAge int `at:"age" json:"-"`
}

type AtBar struct {
	at.Annotation

	at.BaseAnnotation
}

type AtFooBar struct {
	at.Annotation

	AtFoo
	Code int `value:"200" at:"code" json:"-"`
}

type AtFooBaz struct {
	at.Annotation

	AtFoo
	Code int `value:"400" at:"code" json:"-"`
}

type MyObj struct{
	Name string
	Value string
}

type foo struct {
	AtBaz `value:"baz"`
	AtFoo `value:"foo,option 1,option 2" age:"18" fooId:"123"`
	AtBar `value:"bar"`
	AtFooBar `value:"foobar" age:"12"`
	AtFooBaz `value:"foobaz" age:"22"`
	MyObj
}

type bar struct {
	AtFoo `value age:"25"`
	AtBar `value:"bar"`
}

type multipleBar struct {
	AtFoo `value:"foo"`
	AtBar `value:"bar0"`
	Bar1 struct{
		AtBar `value:"bar1"`
		AtBaz `value:"baz"`
	}
	Bar2 struct{
		AtBar `value:"bar2"`
		AtBaz `value:"baz"`
	}
}

type AtIntMap struct {
	at.Annotation

	at.BaseAnnotation
	FieldName string `value:"codes"`
	Codes map[int]string
}

type AtStrMap struct {
	at.Annotation

	at.BaseAnnotation
	FieldName string `value:"messages"`
	Messages map[string]string
}

type foobar struct {
	AtIntMap `200:"success" 404:"not found" 403:"unauthorized"`
	AtStrMap `ok:"successful" failed:"failed"`
}

type AtLevel1 struct {
	at.Annotation

	at.BaseAnnotation
}

type AtLevel2 struct {
	at.Annotation

	AtLevel1
}

type AtLevel3 struct {
	at.Annotation

	AtLevel2
}


type AtLevel4 struct {
	at.Annotation

	AtLevel3
}

type AtLevel5 struct {
	at.Annotation

	AtLevel4
}

type testData struct {
	AtLevel5

	Name string `value:"foo"`
}

type atApiOperation struct {
	at.Operation `value:"testApi" id:"getGreeting" description:"This is the Greeting api for demo"`
}

type AtArray struct {
	at.Annotation

	AtValues []string `at:"values" json:"-"`
}

type fooService struct {
}

func (s *fooService) IHaveAnnotation(_ struct{ AtFoo `value:"bar"` })  {

}

func (s *fooService) IDoNotHaveAnnotation()  {

}

func TestImplementsAnnotation(t *testing.T) {
	log.SetLevel("debug")

	f := new(foo)
	f.Value = "my object value"
	t.Run("test api operation", func(t *testing.T) {
		ao := &atApiOperation{}
		err := annotation.InjectAll(ao)
		assert.Equal(t, nil, err)
	})

	t.Run("test map injection", func(t *testing.T) {
		fb := &foobar{}
		err := annotation.InjectAll(fb)
		assert.Equal(t, nil, err)
	})

	annotations := annotation.GetAnnotations(f)
	t.Run("should check if object contains at.Annotation", func(t *testing.T) {
		assert.NotEqual(t, nil, annotations)

		ok := annotation.Contains(annotations, AtFoo{})
		assert.Equal(t, true, ok)
	})

	t.Run("should find if contains child annotation", func(t *testing.T) {
		ok := annotation.ContainsChild(annotations, AtFoo{})
		assert.Equal(t, true, ok)
	})

	t.Run("should find child annotation", func(t *testing.T) {
		a := annotation.Find(annotations, AtFoo{})
		assert.Equal(t, "AtFoo", a.Field.StructField.Name)
	})

	t.Run("should filter in child annotation", func(t *testing.T) {
		mb := &multipleBar{}
		ma := annotation.GetAnnotations(mb)
		abs := annotation.FilterIn(ma, AtBar{})
		assert.Equal(t, 3, len(abs))
		assert.Equal(t, "AtBar", abs[0].Field.StructField.Name)
		assert.Equal(t, "AtBar", abs[1].Field.StructField.Name)
		assert.Equal(t, "AtBar", abs[2].Field.StructField.Name)
	})

	t.Run("should check if object contains at.Annotation", func(t *testing.T) {
		ok := annotation.Contains(&f.AtFoo, at.Annotation{})
		assert.Equal(t, true, ok)
	})

	t.Run("should get annotation AtFoo", func(t *testing.T) {
		af := annotation.GetAnnotation(f, AtFoo{})
		value, ok := af.Field.StructField.Tag.Lookup("value")
		assert.Equal(t, "foo,option 1,option 2", value)
		assert.Equal(t, true, ok)

		age, ok := af.Field.StructField.Tag.Lookup("age")
		assert.Equal(t, "18", age)
		assert.Equal(t, true, ok)
	})

	t.Run("should report error for invalid type that pass to GetFields", func(t *testing.T) {
		af := annotation.GetAnnotations(123)
		assert.True(t, af == nil)
	})

	t.Run("should report error for invalid type that pass to GetField", func(t *testing.T) {
		a := annotation.GetAnnotation(123, AtFoo{})
		assert.True(t, nil == a)
	})

	t.Run("should inject all annotations", func(t *testing.T) {
		assert.Equal(t, "", f.AtFoo.AtValue)
		assert.Equal(t, 0, f.AtFoo.AtAge)
		err := annotation.InjectAll(f)
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", f.AtFoo.AtValue)
		assert.Equal(t, 18, f.AtFoo.AtAge)

		assert.Equal(t, "bar", f.AtBar.AtValue)
		assert.Equal(t, "my object value", f.Value)
	})

	t.Run("should find annotation AtFoo", func(t *testing.T) {
		as := annotation.FindAll(f, AtFoo{})
		assert.NotEqual(t, 0, len(as))
	})

	t.Run("should find annotation AtBaz", func(t *testing.T) {
		atBazFileds := annotation.FindAll(f, AtBaz{})
		assert.Equal(t, 1, len(atBazFileds))
	})

	t.Run("should notify bad syntax for struct tag pair", func(t *testing.T) {
		// notify bad syntax for struct tag pair
		b := new(bar)
		err := annotation.InjectAll(b)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, "bad syntax for struct tag pair", err.Error())
	})

	t.Run("should inject to object", func(t *testing.T) {
		fo := foo{}
		err := annotation.InjectAll(fo)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should check if an object implements Annotation", func(t *testing.T) {
		ok := annotation.Contains(f, AtFoo{})
		assert.Equal(t, true, ok)
	})

	t.Run("should inject annotation into sub struct", func(t *testing.T) {
		var fb struct{at.GetMapping `value:"/path/to/api" foo:"bar"`}
		err := annotation.InjectAll(&fb)
		assert.Equal(t, nil, err)
	})

	t.Run("should report error when inject nil object", func(t *testing.T) {
		err := annotation.InjectAll(nil)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should report error when inject invalid tagged annotation", func(t *testing.T) {
		ff := annotation.GetAnnotation(&bar{}, AtFoo{})
		err := annotation.Inject(ff)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should get annotation by type", func(t *testing.T) {
		var fb struct{at.PostMapping `value:"/path/to/api"`}
		f := annotation.GetAnnotation(reflect.TypeOf(&fb), at.PostMapping{})
		assert.NotEqual(t, nil, f)
		assert.Equal(t, "PostMapping", f.Field.StructField.Name)
		assert.Equal(t, false, f.Field.Value.IsValid())
	})

	t.Run("should inject into a field", func(t *testing.T) {
		ff := new(foo)
		ann := annotation.GetAnnotation(ff, AtBaz{})
		assert.NotEqual(t, nil, ann)
		err := annotation.Inject(ann)
		assert.Equal(t, nil, err)
	})

	t.Run("should get and inject into multiple sub annotations", func(t *testing.T) {
		ma := &multipleBar{}
		maf := annotation.GetAnnotations(ma)
		assert.Equal(t, 2, len(maf.Items))
		assert.Equal(t, 2, len(maf.Children))

		err := annotation.InjectAll(ma)
		assert.Equal(t, nil, err)
		assert.Equal(t, "bar1", ma.Bar1.AtBar.AtValue)
		assert.Equal(t, "bar2", ma.Bar2.AtBar.AtValue)
		assert.Equal(t, "baz", ma.Bar1.AtBaz.AtValue)

		fa := annotation.Find(maf, AtFoo{})
		assert.Equal(t, "foo", fa.Field.Value.Interface().(AtFoo).AtValue)
	})

	t.Run("should get from nil interface", func(t *testing.T) {
		a := annotation.GetAnnotation((*multipleBar)(nil), AtFoo{})
		assert.Equal(t, false, a.Field.Value.IsValid())
	})

	t.Run("should get from nil interface", func(t *testing.T) {
		a := annotation.GetAnnotation((*annotation.Annotations)(nil), AtFoo{})
		assert.True(t, nil == a)
	})

	t.Run("should get nil from nil", func(t *testing.T) {
		a := annotation.GetAnnotation(nil, AtFoo{})
		assert.True(t, nil == a)
	})
	t.Run("should get nil from nil", func(t *testing.T) {
		a := annotation.GetAnnotations(nil)
		assert.True(t, nil == a)
	})

	t.Run("should get annotation from function return object", func(t *testing.T) {
		a := annotation.GetAnnotation(func() *multipleBar {
			return &multipleBar{}
		}, AtFoo{})
		assert.Equal(t, "AtFoo", a.Field.StructField.Name)
	})

	t.Run("should find each annotation deeply", func(t *testing.T) {
		td := &testData{}

		assert.Equal(t, false, annotation.IsAnnotation(nil))

		assert.Equal(t, false, annotation.IsAnnotation(123))

		assert.Equal(t, false, annotation.IsAnnotation("abc"))

		assert.Equal(t, false, annotation.IsAnnotation(true))

		assert.Equal(t, false, annotation.IsAnnotation([]string{"abc"}))

		assert.Equal(t, false, annotation.IsAnnotation(td))

		l5 := annotation.GetAnnotation(td, AtLevel5{})
		assert.Equal(t, "AtLevel5", l5.Field.StructField.Name)
		assert.Equal(t, true, annotation.IsAnnotation(l5.Field.Value))

		l4 := annotation.GetAnnotation(td, AtLevel4{})
		assert.Equal(t, "AtLevel4", l4.Field.StructField.Name)
		assert.Equal(t, true, annotation.IsAnnotation(l4.Field.Value))

		l3 := annotation.GetAnnotation(td, AtLevel3{})
		assert.Equal(t, "AtLevel3", l3.Field.StructField.Name)
		assert.Equal(t, true, annotation.IsAnnotation(l3.Field.Value))

		l2 := annotation.GetAnnotation(td, AtLevel2{})
		assert.Equal(t, "AtLevel2", l2.Field.StructField.Name)
		assert.Equal(t, true, annotation.IsAnnotation(l2.Field.Value))

		l1 := annotation.GetAnnotation(td, AtLevel1{})
		assert.Equal(t, "AtLevel1", l1.Field.StructField.Name)
		assert.Equal(t, true, annotation.IsAnnotation(l1.Field.Value))

		base := annotation.GetAnnotation(td, at.BaseAnnotation{})
		assert.Equal(t, "BaseAnnotation", base.Field.StructField.Name)
		assert.Equal(t, true, annotation.IsAnnotation(base.Field.Value))

		ann := annotation.GetAnnotation(td, at.Annotation{})
		assert.Equal(t, "Annotation", ann.Field.StructField.Name)
		assert.Equal(t, true, annotation.IsAnnotation(ann.Field.Value))
	})

	t.Run("should inject into annotation", func(t *testing.T) {
		type foo struct {
			at.Schema `value:"array" type:"string" description:"This is a test parameter"`

			Name string `json:"name"`
		}
		f := &foo{Name: "foo"}
		a := annotation.GetAnnotation(f, at.Schema{})
		err := annotation.Inject(a)
		assert.Equal(t, nil, err)
		ao := a.Field.Value.Interface().(at.Schema)
		assert.Equal(t, "array", ao.AtValue)
		assert.Equal(t, "string", ao.AtType)
		assert.Equal(t, "This is a test parameter", ao.AtDescription)
	})

	t.Run("should inject into annotation", func(t *testing.T) {
		type foo struct {
			at.Parameter `value:"foo" in:"path" type:"integer" description:"This is a test parameter"`
		}
		f := &foo{}
		a := annotation.GetAnnotation(f, at.Parameter{})
		err := annotation.Inject(a)
		assert.Equal(t, nil, err)
		ao := a.Field.Value.Interface().(at.Parameter)
		assert.Equal(t, "integer", ao.AtType)
		assert.Equal(t, "path", ao.AtIn)
		assert.Equal(t, "This is a test parameter", ao.AtDescription)
	})

	t.Run("should parse complicate annotations", func(t *testing.T) {
		atTest := &struct {
			at.PostMapping `value:"/"`
			at.Operation   `id:"Create Employee" description:"This is the employee creation api"`
			at.Consumes    `values:"application/json"`
			at.Produces    `values:"application/json"`
			Parameters     struct {
				at.Parameter `name:"token" in:"header" type:"string" description:"JWT token (fake token - for demo only)" `
				Body struct {
					at.Parameter `name:"employee" in:"body" description:"Employee request body" `
					foo
				}
			}
			Responses struct {
				StatusOK struct {
					at.Response `code:"200" description:"returns a employee with ID"`
					XRateLimit struct {
						at.Header `value:"X-Rate-Limit" type:"integer" format:"int32" description:"calls per hour allowed by the user"`
					}
					XExpiresAfter struct{
						at.Header `value:"X-Expires-After" type:"string" format:"date-time" description:"date in UTC when token expires"`
					}
					bar
				}
			}
		}{}

		ans := annotation.GetAnnotations(atTest)
		err := annotation.InjectAll(atTest)
		assert.Equal(t, nil, err)
		assert.Equal(t, "Header", ans.Children[1].Children[0].Children[0].Items[0].Field.StructField.Name)
		assert.Equal(t, "X-Rate-Limit", atTest.Responses.StatusOK.XRateLimit.AtValue)
	})

	t.Run("should inject array", func(t *testing.T) {
		type testArray struct {
			AtArray `values:"foo,bar,baz"`
		}
		ta := &testArray{}
		err := annotation.InjectAll(ta)
		assert.Equal(t, nil, err)
		assert.Equal(t, 3, len(ta.AtValues))
	})

	t.Run("find annotated methods", func(t *testing.T) {
		fs := new(fooService)
		methods, annotations := annotation.FindAnnotatedMethods(fs, AtFoo{})
		assert.Equal(t, 1, len(methods))
		assert.Equal(t, 1, len(annotations))
	})
}

