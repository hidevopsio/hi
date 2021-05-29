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

package inject

import (
	"errors"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
	"strings"
)

const (
	value          = "value"
)

var (
	// ErrNotImplemented the interface is not implemented
	ErrNotImplemented = errors.New("[inject] interface is not implemented")

	// ErrInvalidObject the object is invalid
	ErrInvalidObject = errors.New("[inject] invalid object")

	// ErrInvalidTagName the tag name is invalid
	ErrInvalidTagName = errors.New("[inject] invalid tag name, e.g. exampleTag")

	// ErrSystemConfiguration system is not configured
	ErrSystemConfiguration = errors.New("[inject] system is not configured")

	// ErrInvalidFunc the function is invalid
	ErrInvalidFunc = errors.New("[inject] invalid func")

	// ErrInvalidMethod the function is invalid
	ErrInvalidMethod = errors.New("[inject] invalid method")

	// ErrFactoryIsNil factory is invalid
	ErrFactoryIsNil = errors.New("[inject] factory is nil")

	ErrAnnotationsIsNil = fmt.Errorf("err: annotations is nil")

	tagsContainer []Tag

	//instancesMap cmap.ConcurrentMap
	//appFactory factory.ConfigurableFactory
)

// Inject is the interface for inject tag
type Inject interface {
	DefaultValue(object interface{}) error
	IntoObject(object interface{}) error
	IntoObjectValue(object reflect.Value, property string, tags ...Tag) error
	IntoMethod(object interface{}, m interface{}) (retVal interface{}, err error)
	IntoFunc(object interface{}) (retVal interface{}, err error)
	IntoAnnotations(annotations *annotation.Annotations) (err error)
}

type inject struct {
	factory factory.InstantiateFactory
}

// NewInject is the constructor of inject
func NewInject(factory factory.InstantiateFactory) Inject {
	return &inject{factory: factory}
}

// InitTag init tag implements
func InitTag(tag Tag) (t Tag) {
	if annotation.Contains(tag, at.Tag{}) {
		err := annotation.InjectAll(tag)
		if err == nil {
			t = tag
		}
	}
	return
}

// AddTag add new tag
func AddTag(tag Tag) {
	if tag != nil {
		t := InitTag(tag)
		if t != nil {
			tagsContainer = append(tagsContainer, t)
		}
	}
}

func (i *inject) getInstance(typ reflect.Type) (inst interface{}) {
	n := reflector.GetLowerCamelFullNameByType(typ)
	inst = i.factory.GetInstance(n)
	return
}

// DefaultValue injects instance into the tagged field with `inject:"instanceName"`
func (i *inject) DefaultValue(object interface{}) error {
	return i.IntoObjectValue(reflect.ValueOf(object), "", InitTag(new(defaultTag)))
}

func (i *inject) IntoAnnotations(annotations *annotation.Annotations) (err error) {
	if annotations == nil {
		err = ErrAnnotationsIsNil
		return
	}

	// inject annotation
	for _, a := range annotations.Items {
		err = annotation.Inject(a)
		if err == nil && a.Field.Value.IsValid() {
			err = i.IntoObjectValue(a.Field.Value.Addr(), "")
		}
	}

	for _, c := range annotations.Children {
		err = i.IntoAnnotations(c)
	}
	return
}

// IntoObject injects instance into the tagged field with `inject:"instanceName"`
func (i *inject) IntoObject(object interface{}) (err error) {
	//
	//err = annotation.InjectAll(object)
	//if err != nil {
	//	log.Debug(err)
	//}

	// inject into value
	err = i.IntoObjectValue(reflect.ValueOf(object), "")

	// inject into annotations
	if err == nil {
		annotations := annotation.GetAnnotations(object)
		err = i.IntoAnnotations(annotations)
	}
	return
}

func (i *inject) convert(f reflect.StructField, src interface{}) (fov reflect.Value) {
	fov = reflect.ValueOf(src)

	// convert slice
	switch src.(type) {
	//case []string:
	case []interface{}:
		switch f.Type.Elem().Kind() {
		case reflect.String:
			var sv []string
			src := src.([]interface{})
			for _, elm := range src {
				sv = append(sv, elm.(string))
			}
			fov = reflect.ValueOf(sv)
		}
	}
	//log.Debugf("Injected slice %v.(%v) into %v.%v", src, fov.Type(), fov.Type(), f.Name)
	return
}

// IntoObjectValue injects instance into the tagged field with `inject:"instanceName"`
func (i *inject) IntoObjectValue(object reflect.Value, property string, tags ...Tag) error {
	var err error

	//// TODO refactor IntoObject
	//if appFactory == nil {
	//	return ErrSystemConfiguration
	//}

	obj := reflector.Indirect(object)
	if obj.Kind() != reflect.Struct {
		log.Debugf("[inject] ignore object: %v, kind: %v", object, obj.Kind())
		return ErrInvalidObject
	}

	var targetTags []Tag
	if len(tags) != 0 {
		targetTags = tags
	} else {
		targetTags = tagsContainer
	}

	// field injection
	for _, f := range reflector.DeepFields(object.Type()) {
		var injectedObject interface{}
		var prop string
		pn, ok := f.Tag.Lookup("json")
		if ok {
			pns := strings.Split(pn, ",")
			if len(pns) > 1 {
				prop = pns[0]
			}
		} else {
			prop = str.ToLowerCamel(f.Name)
		}

		// debug prints
		//n := reflector.GetLowerCamelFullNameByType(f.Type)
		//log.Debugf("%v : %v",property + "." + prop, n)

		if property != "" {
			prop = property + "." + prop
		}

		//log.Debugf("parent: %v, name: %v, type: %v, tag: %v", obj.Type(), f.Name, f.Type, f.Tag)
		// check if object has value field to be injected

		ft := f.Type
		if f.Type.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		// set field object
		var fieldObjValue reflect.Value
		if obj.IsValid() && obj.Kind() == reflect.Struct {
			// TODO: consider embedded property
			//log.Debugf("inject: %v.%v", obj.Type(), f.Name)
			fieldObjValue = obj.FieldByName(f.Name)
		}

		// TODO: assume that the f.Name of value and inject tag is not the same
		injectedObject = i.getInstance(f.Type)
		if injectedObject == nil {
			for _, tagImpl := range targetTags {
				tagImpl.Init(i.factory)
				injectedObject = tagImpl.Decode(object, f, prop)
				if injectedObject != nil {
					break
				}
			}
		}

		// assign value to struct field
		if ft.Kind() != reflect.Struct || annotation.Contains(injectedObject, at.AutoWired{}) {
			if injectedObject != nil && fieldObjValue.CanSet() {
				fov := i.convert(f, injectedObject)
				if fov.Type().AssignableTo(fieldObjValue.Type()) {
					fieldObjValue.Set(fov)
					//} else {
					//	log.Errorf("unmatched type %v against %v", fov.Type(), fieldObj.Type())
				}
			}
		}

		//log.Debugf("- kind: %v, %v, %v, %v", obj.Kind(), object.Type(), fieldObj.Type(), f.Name)
		//log.Debugf("isValid: %v, canSet: %v", fieldObj.IsValid(), fieldObj.CanSet())
		filedObject := reflect.Indirect(fieldObjValue)
		filedKind := filedObject.Kind()
		canNested := filedKind == reflect.Struct
		if canNested && fieldObjValue.IsValid() && fieldObjValue.CanSet() && filedObject.Type() != obj.Type() {
			err = i.IntoObjectValue(fieldObjValue, prop, tags...)
		}
	}

	//inject property set
	//if atFields := annotation.FindAll(object, at.ConfigurationProperties{}); len(atFields) > 0 {
	//	obj := object.Interface()
	//	err = i.factory.Builder().Load(obj)
	//}
	return err
}

func (i *inject) parseFuncOrMethodInput(inType reflect.Type) (paramValue reflect.Value, ok bool) {
	inType = reflector.IndirectType(inType)
	inst := i.getInstance(inType)
	ok = true
	if inst == nil {
		//log.Debug(inType.Kind())
		switch inType.Kind() {
		// interface and slice creation is not supported
		case reflect.Interface, reflect.Slice:
			ok = false
			break
		default:
			// should find instance in the component container first

			// if it is not found, then create new instance
			paramValue = reflect.New(inType)
			if annotation.IsAnnotation(inType) {
				inst = paramValue.Elem().Interface()
			} else {
				inst = paramValue.Interface()
			}
			// TODO: inTypeName
			i.factory.SetInstance(inst)
		}
	}

	if inst != nil {
		paramValue = reflect.ValueOf(inst)
	}
	return
}

// IntoFunc inject object into func and return instance
func (i *inject) IntoFunc(object interface{}) (retVal interface{}, err error) {
	fn := reflect.ValueOf(object)
	if fn.Kind() == reflect.Func {
		numIn := fn.Type().NumIn()
		inputs := make([]reflect.Value, numIn)
		// TODO: should load function inputs when resolving dependencies to improve performance
		for n := 0; n < numIn; n++ {
			fnInType := fn.Type().In(n)
			//expectedTypName := reflector.GetLowerCamelFullNameByType(fnInType)
			//log.Debugf("expected: %v", expectedTypName)
			val, ok := i.parseFuncOrMethodInput(fnInType)
			if ok {

				inputs[n] = val
				//log.Debugf("Injected %v into func parameter %v", val, fnInType)
			} else {
				return nil, fmt.Errorf("%v is not injected", fnInType.Name())
			}

			paramValue := reflect.Indirect(val)
			if val.IsValid() && paramValue.IsValid() && paramValue.Kind() == reflect.Struct {
				err = i.IntoObject(val.Interface())
			}
		}
		results := fn.Call(inputs)
		if len(results) != 0 {
			retVal = results[0].Interface()
			return
		}
		return
	}
	err = ErrInvalidFunc
	return
}

//IntoMethod inject object into func and return instance
//TODO: IntoMethod or IntoFunc should accept metaData, because it contains dependencies
func (i *inject) IntoMethod(object interface{}, m interface{}) (retVal interface{}, err error) {
	if object != nil && m != nil {
		switch m.(type) {
		case reflect.Method:
			method := m.(reflect.Method)
			numIn := method.Type.NumIn()
			inputs := make([]reflect.Value, numIn)
			inputs[0] = reflect.ValueOf(object)
			var ann interface{}
			for n := 1; n < numIn; n++ {
				fnInType := method.Type.In(n)
				if annotation.IsAnnotation(fnInType) {
					ann = fnInType
				}
				val, ok := i.parseFuncOrMethodInput(fnInType)
				if ok {
					inputs[n] = val
				} else {
					if reflect.TypeOf(at.AllowNil{}) == ann || annotation.Contains(ann, at.AllowNil{}) {
						inputs[n] = reflect.Zero(fnInType)
					} else {
						return nil, fmt.Errorf("%v is not injected", fnInType.Name())
					}
				}

				paramObject := reflect.Indirect(val)
				if val.IsValid() && paramObject.IsValid() && paramObject.Kind() == reflect.Struct {
					err = i.IntoObject(val.Interface())
				}
			}
			results := method.Func.Call(inputs)
			if len(results) != 0 {
				retVal = results[0].Interface()
				return
			}
		}
	}
	err = ErrInvalidMethod
	return
}
