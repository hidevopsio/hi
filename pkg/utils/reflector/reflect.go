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

// Package reflector provides utilities for reflection
package reflector

import (
	"errors"
	"reflect"
	"runtime"
	"strings"

	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
)

var (
	// ErrInvalidInput means that the input is invalid
	ErrInvalidInput = errors.New("input is invalid")

	// ErrInvalidMethod means that the method is invalid
	ErrInvalidMethod = errors.New("method is invalid")

	// ErrInvalidFunc means that the func is invalid
	ErrInvalidFunc = errors.New("func is invalid")

	// ErrFieldCanNotBeSet means that the field can not be set
	ErrFieldCanNotBeSet = errors.New("field can not be set")
)

// NewReflectType create instance by tyep
func NewReflectType(st interface{}) interface{} {
	ct := reflect.TypeOf(st)
	co := reflect.New(ct)
	cp := co.Elem().Addr().Interface()
	return cp
}

// Validate validate value
func Validate(toValue interface{}) (*reflect.Value, error) {

	to := Indirect(reflect.ValueOf(toValue))

	// Return is from value is invalid
	if !to.IsValid() {
		return nil, errors.New("value is not valid")
	}

	if !to.CanAddr() {
		return nil, errors.New("value is unaddressable")
	}

	return &to, nil
}

// DeepFields iterate struct field
func DeepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = IndirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				vk :=  IndirectType(v.Type).Kind()
				if vk != reflect.Struct && vk != reflect.Interface {
					fields = append(fields, v)
				} else {
					fields = append(fields, DeepFields(v.Type)...)
				}
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

// Indirect get indirect value
func Indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

// IndirectValue get indirect value
func IndirectValue(object interface{}) (reflectValue reflect.Value) {
	switch object.(type) {
	case reflect.Value:
		reflectValue = object.(reflect.Value)
	case reflect.Type:
		return
	default:
		reflectValue = reflect.ValueOf(object)
	}
	reflectValue = Indirect(reflectValue)
	return
}

// IndirectType get indirect type
func IndirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

// GetFieldValue get field value
func GetFieldValue(f interface{}, name string) reflect.Value {
	r := reflect.ValueOf(f)
	fv := reflect.Indirect(r).FieldByName(name)

	return fv
}


// FindFieldByTag find field by tag
func FindFieldByTag(obj interface{}, key, name string) (field reflect.StructField, ok bool) {
	typ, ok := GetObjectType(obj)
	if ok {
		for _, f := range DeepFields(typ) {
			if f.Tag.Get(key) == name {
				field = f
				ok = true
				break
			}
		}
	}

	return
}


// SetFieldValue set field value
func SetFieldValue(object interface{}, name string, value interface{}) error {

	obj := Indirect(reflect.ValueOf(object))

	if !obj.IsValid() {
		return ErrInvalidInput
	}

	if obj.Kind() != reflect.Struct {
		return ErrInvalidInput
	}

	fieldObj := obj.FieldByName(name)

	if !fieldObj.CanSet() {
		return ErrFieldCanNotBeSet
	}

	fov := reflect.ValueOf(value)
	fieldObj.Set(fov)

	//log.Debugf("Set %v.(%v) into %v.%v", value, fov.Type(), obj.Type(), name)
	return nil
}

// GetKind get kind
func GetKind(kind reflect.Kind) reflect.Kind {

	// Check each condition until a case is true.
	switch {

	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int

	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint

	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32

	default:
		return kind
	}
}

// GetKindByValue get kind by value
func GetKindByValue(val reflect.Value) reflect.Kind {
	return GetKind(val.Kind())
}

// GetKindByType get kind by type
func GetKindByType(typ reflect.Type) reflect.Kind {
	return GetKind(typ.Kind())
}

// ValidateReflectType validate reflect type and iterate all fields
func ValidateReflectType(obj interface{}, callback func(value *reflect.Value, reflectType reflect.Type, fieldSize int, isSlice bool) error) error {
	v, err := Validate(obj)
	if err != nil {
		return err
	}

	t := IndirectType(v.Type())

	isSlice := false
	fieldSize := 1
	if v.Kind() == reflect.Slice {
		isSlice = true
		fieldSize = v.Len()
	}

	if callback != nil {
		return callback(v, t, fieldSize, isSlice)
	}

	return err
}

// GetType get object data type
func GetType(data interface{}) (typ reflect.Type, err error) {
	dv := Indirect(reflect.ValueOf(data))

	// Return is from value is invalid
	if !dv.IsValid() {
		err = ErrInvalidInput
		return
	}
	typ = dv.Type()

	//log.Debugf("%v %v %v %v %v", dv, typ, typ.String(), typ.Name(), typ.PkgPath())
	return
}

func parseFuncName(name string) string {
	n := strings.LastIndexByte(name, byte('.'))
	if n > 0 {
		name = name[n+1:]
		n = strings.LastIndexByte(name, byte('-'))
		if n > 0 {
			name = name[:n]
		}
		n = strings.LastIndexByte(name, byte(')'))
		if n > 0 {
			name = name[:n]
		}
	}
	return name
}

// GetFuncName get func name
func GetFuncName(fn interface{}) (name string) {
	val := reflect.ValueOf(fn)
	kind := val.Kind()
	if kind == reflect.Func {
		name = runtime.FuncForPC(val.Pointer()).Name()
		name = parseFuncName(name)
	}
	return
}

// GetName get object name
func GetName(data interface{}) (name string) {

	typ, err := GetType(data)
	if err == nil {
		name = typ.Name()
	}
	return
}

// GetLowerCamelName get lower case object name
func GetLowerCamelName(object interface{}) (name string) {
	name = GetName(object)
	name = str.ToLowerCamel(name)
	return name
}

// HasField check if has specific field
func HasField(object interface{}, name string) bool {
	r := reflect.ValueOf(object)
	fv := reflect.Indirect(r).FieldByName(name)

	return fv.IsValid()
}


// GetMethodsByAnnotation call method
func GetMethodsByAnnotation(object interface{}, name string, args ...interface{}) (interface{}, error) {
	objVal := reflect.ValueOf(object)
	method, ok := objVal.Type().MethodByName(name)
	if ok {
		numIn := method.Type.NumIn()
		inputs := make([]reflect.Value, numIn)
		inputs[0] = objVal
		for i, arg := range args {
			inputs[i+1] = reflect.ValueOf(arg)
		}
		results := method.Func.Call(inputs)
		if len(results) != 0 {
			return results[0].Interface(), nil
		}
		return nil, nil
	}
	return nil, ErrInvalidMethod
}

// CallMethodByName call method
func CallMethodByName(object interface{}, name string, args ...interface{}) (interface{}, error) {
	objVal := reflect.ValueOf(object)
	method, ok := objVal.Type().MethodByName(name)
	if ok {
		numIn := method.Type.NumIn()
		inputs := make([]reflect.Value, numIn)
		inputs[0] = objVal
		for i, arg := range args {
			inputs[i+1] = reflect.ValueOf(arg)
		}
		results := method.Func.Call(inputs)
		if len(results) != 0 {
			return results[0].Interface(), nil
		}
		return nil, nil
	}
	return nil, ErrInvalidMethod
}

// CallFunc call function
func CallFunc(object interface{}, args ...interface{}) (interface{}, error) {
	fn := reflect.ValueOf(object)
	if fn.Kind() == reflect.Func {
		numIn := fn.Type().NumIn()
		inputs := make([]reflect.Value, numIn)
		for i, arg := range args {
			inputs[i] = reflect.ValueOf(arg)
		}
		results := fn.Call(inputs)
		if len(results) != 0 {
			return results[0].Interface(), nil
		}
		return nil, nil
	}
	return nil, ErrInvalidFunc
}

// HasEmbeddedField check if has embedded fieled
func HasEmbeddedField(object interface{}, name string) bool {
	//log.Debugf("HasEmbeddedField: %v", name)
	typ := IndirectType(reflect.TypeOf(object))
	if typ.Kind() != reflect.Struct {
		return false
	}
	field, ok := typ.FieldByName(name)
	return field.Anonymous && ok
}

// GetEmbeddedFieldType check if has embedded fieled
func GetEmbeddedFieldType(object interface{}, expected interface{}) (field reflect.StructField, ok bool) {
	if object == nil {
		return
	}
	typ, got := GetObjectType(object)
	if got {
		expectedTyp, _ := GetObjectType(expected)

		field, ok = GetEmbeddedFieldByType(typ, expectedTyp.Name(), expectedTyp.Kind())
	}

	return
}

// HasEmbeddedFieldType check if has embedded fieled
func HasEmbeddedFieldType(object interface{}, expected interface{}) (found bool) {
	_, found = GetEmbeddedFieldType(object, expected)
	return
}

// GetEmbeddedFieldByType get embedded interface field by type
func GetEmbeddedFieldByType(typ reflect.Type, nameOrType interface{}, kind ...reflect.Kind) (field reflect.StructField, ok bool) {
	expectedKind := reflect.Interface
	if len(kind) > 0 {
		expectedKind = kind[0]
	}
	var name string
	switch nameOrType.(type) {
	case string:
		name = nameOrType.(string)
	default:
		name = GetName(nameOrType)
	}

	k := typ.Kind()
	if k == reflect.Struct {
		numField := typ.NumField()
		for i := 0; i < numField; i++ {
			v := typ.Field(i)
			if v.Anonymous {
				if v.Type.Kind() == expectedKind && (name == "" || v.Name == name) {
					field = v
					ok = true
					break
				} else {
					field, ok = GetEmbeddedFieldByType(v.Type, name, kind...)
					if ok {
						break
					}
				}
			}
		}
	}
	return
}

// GetEmbeddedFieldsByType get embedded interface fields by type
func GetEmbeddedFieldsByType(typ reflect.Type, kind ...reflect.Kind) (fields []reflect.StructField) {
	expectedKind := reflect.Interface
	if len(kind) > 0 {
		expectedKind = kind[0]
	}
	if typ == nil {
		return
	}
	k := typ.Kind()
	if k == reflect.Struct {
		numField := typ.NumField()
		for i := 0; i < numField; i++ {
			v := typ.Field(i)
			if v.Anonymous {
				if v.Type.Kind() == expectedKind {
					fields = append(fields, v)
				} else {
					f := GetEmbeddedFieldsByType(v.Type)
					fields = append(fields, f...)
				}
			}
		}
	}
	return
}

// GetEmbeddedFields get embedded interface fields
func GetEmbeddedFields(object interface{}, kind ...reflect.Kind) (fields []reflect.StructField) {
	if object == nil {
		return
	}
	typ, _ := GetObjectType(object)
	fields = GetEmbeddedFieldsByType(typ, kind...)
	return
}

// GetEmbeddedField get embedded interface field
func GetEmbeddedField(object interface{}, nameOrType interface{}, dataTypes ...reflect.Kind) (field reflect.StructField) {
	if object == nil {
		return
	}

	typ, ok := GetObjectType(object)
	if ok {
		f, ok := GetEmbeddedFieldByType(typ, nameOrType, dataTypes...)
		if ok {
			field = f
		}
	}
	return
}

// FindEmbeddedFieldTag find embedded field tag
func FindEmbeddedFieldTag(object interface{}, fieldName, tagName string) (t string, ok bool) {
	f := GetEmbeddedField(object, fieldName)
	t, ok = f.Tag.Lookup(tagName)
	return
}

// ParseObjectName e.g. ExampleObject => example
func ParseObjectName(obj interface{}, eliminator string) string {
	name := GetName(obj)
	name = strings.Replace(name, eliminator, "", -1)
	name = str.LowerFirst(name)
	return name
}

// ParseObjectPkgName e.g. ExampleObject => example
func ParseObjectPkgName(obj interface{}) string {

	typ := IndirectType(reflect.TypeOf(obj))
	name := io.DirName(typ.PkgPath())

	return name
}

// GetPkgPath get the package patch
func GetPkgPath(object interface{}) string {
	objType := IndirectType(reflect.TypeOf(object))
	return objType.PkgPath()
}

// GetObjectType get the function output data type
func GetObjectType(object interface{}) (typ reflect.Type, ok bool) {
	if object == nil {
		return nil, false
	}
	var reflectType reflect.Type
	var reflectValue reflect.Value
	switch object.(type) {
	case reflect.Value:
		reflectValue = object.(reflect.Value)
		reflectType = reflectValue.Type()
	case reflect.Type:
		reflectType = object.(reflect.Type)
	default:
		reflectValue = reflect.ValueOf(object)
		reflectType = reflectValue.Type()
	}
	typName := reflectType.Name()
	typKind := reflectType.Kind()
	//log.Debugf("type: %v type name: %v, kind: %v", t, typName, typKind)
	if typKind == reflect.Func {
		numOut := reflectType.NumOut()
		if numOut > 0 {
			typ = IndirectType(reflectType.Out(0))
			ok = true
		}
	} else if typKind == reflect.Struct && typName == "Method" {
		method := object.(reflect.Method)
		methodTyp := method.Func.Type()
		numOut := methodTyp.NumOut()
		if numOut > 0 {
			typ = IndirectType(methodTyp.Out(0))
			ok = true
		}
	} else {
		// TODO: check if it effects others
		typ = IndirectType(reflectType)
		ok = true
	}

	return
}

// GetFullName get the object name with package name, e.g. pkg.Object
func GetFullName(object interface{}) (name string) {

	pn, n := GetPkgAndName(object)
	name = pn + "." + n

	return
}

// GetLowerCamelFullName get the object name with package name, e.g. pkg.objectName
func GetLowerCamelFullName(object interface{}) (name string) {

	pn, n := GetPkgAndName(object)
	name = pn + "." + str.ToLowerCamel(n)

	return
}

// GetPkgAndName get the package name and the object name with, e.g. pkg, Object
func GetPkgAndName(object interface{}) (pkgName, name string) {

	typ, ok := GetObjectType(object)
	if ok {
		pkgName = io.DirName(typ.PkgPath())
		name = typ.Name()
	}
	return
}

// GetLowerCamelFullNameByType get the object name with package name by type, e.g. pkg.Object
func GetLowerCamelFullNameByType(objType reflect.Type) (name string) {
	indTyp := IndirectType(objType)
	depPkgName := io.DirName(indTyp.PkgPath())
	name = depPkgName + "." + str.ToLowerCamel(indTyp.Name())
	return
}

// IsValidObjectType check if is valid object type
func IsValidObjectType(inst interface{}) bool {
	val := reflect.ValueOf(inst)
	//log.Println(val.Kind())
	//log.Println(reflect.Indirect(val).Kind())
	if val.Kind() == reflect.Ptr && reflect.Indirect(val).Kind() == reflect.Struct {
		return true
	}
	return false
}

// Implements
func Implements(object interface{}, interfaceType interface{}) (ok bool) {
	s := reflect.ValueOf(object)
	typ := s.Type()
	modelType := reflect.TypeOf(interfaceType).Elem()
	ok = typ.Implements(modelType)

	return
}
