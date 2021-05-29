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

// Notes: this source code is originally copied from https://github.com/jinzhu/copier
// The MIT License (MIT)
//
// Copyright (c) 2015 Jinzhu
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package copier provides utility that copy element between structs
package copier

import (
	"database/sql"
	"errors"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
)

// Config configurations for Copy function
type Config struct {
	IgnoreEmptyValue        bool
}

// IgnoreEmptyValue option to config IgnoreEmptyValue, any empty or nil value will not copy from source to destination
func IgnoreEmptyValue(config *Config) {
	config.IgnoreEmptyValue = true
}

func set(to, from reflect.Value, config *Config) bool {
	if from.IsValid() && to.IsValid() {
		if to.Kind() == reflect.Ptr {
			isFromNil := false
			if from.Kind() == reflect.Ptr && from.IsNil() {
				isFromNil = true
			}

			if config.IgnoreEmptyValue && isFromNil {
				return true
			}
			//set `to` to nil if from is nil
			if isFromNil {
				to.Set(reflect.Zero(to.Type()))
				return true
			}

			if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			fVal := from.Convert(to.Type())
			fiVal := fVal.Interface()
			switch fiVal.(type) {
			case string:
				if config.IgnoreEmptyValue && fiVal == "" {
					return true
				}
			}
			if to.Kind() == reflect.Slice {
				if from.Kind() == reflect.Slice {
					if from.Len() == 0 {
						return true
					}
				}
			}

			to.Set(fVal)
		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			err := scanner.Scan(from.Interface())
			if err != nil {
				return false
			}
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem(), config)
		}
		return false
	}
	return false
}

func copy(toValue interface{}, fromValue interface{}, config *Config) (err error) {
	var (
		isSlice bool
		amount  = 1
		from    = reflector.Indirect(reflect.ValueOf(fromValue))
		to      = reflector.Indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}
	// Just set it if possible to assign
	if from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return
	}

	fromType := reflector.IndirectType(from.Type())
	toType := reflector.IndirectType(to.Type())

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = reflector.Indirect(from.Index(i))
			} else {
				source = reflector.Indirect(from)
			}

			// dest
			dest = reflector.Indirect(reflect.New(toType).Elem())
		} else {
			source = reflector.Indirect(from)
			dest = reflector.Indirect(to)
		}

		// Copy from field to field or method
		for _, field := range reflector.DeepFields(fromType) {
			name := field.Name

			if fromField := source.FieldByName(name); fromField.IsValid() {
				// has field
				if toField := dest.FieldByName(name); toField.IsValid() {
					if toField.CanSet() {
						if !set(toField, fromField, config) {
							err = copy(toField.Addr().Interface(), fromField.Interface(), config)
						}
					}
				} else {
					// try to set to method
					var toMethod reflect.Value
					toMethod = dest.MethodByName(name)
					if dest.CanAddr() {
						toMethod = dest.Addr().MethodByName(name)
					}

					if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
						toMethod.Call([]reflect.Value{fromField})
					}
				}
			}
		}

		// Copy from method to field
		for _, field := range reflector.DeepFields(toType) {
			name := field.Name

			var fromMethod reflect.Value
			if source.CanAddr() {
				fromMethod = source.Addr().MethodByName(name)
			} else {
				fromMethod = source.MethodByName(name)
			}

			if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 {
				if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
					values := fromMethod.Call([]reflect.Value{})
					if len(values) >= 1 {
						set(toField, values[0], config)
					}
				}
			}
		}

		if isSlice {
			var fromSliceVal reflect.Value
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				fromSliceVal = reflect.Append(to, dest.Addr())
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				fromSliceVal = reflect.Append(to, dest)
			}
			if !(config.IgnoreEmptyValue && fromSliceVal.IsNil()) {
				to.Set(fromSliceVal)
			}
		}
	}
	return
}

// Copy: copy things from source to destination, Use github.com/jinzhu/copier instead
// Deprecated
func Copy(toValue interface{}, fromValue interface{}, opts ...func(*Config)) (err error) {
	config := &Config{}

	for _, opt := range opts {
		opt(config)
	}

	return copy(toValue, fromValue, config)
}

func copyMap(dst, src map[string]interface{}, config *Config) {
	for k, v := range src {
		dv := dst[k]
		if config.IgnoreEmptyValue && v == nil {
			continue
		}
		switch v.(type) {
		case map[string]interface{}:
			var dm map[string]interface{}
			if dv == nil {
				dm = make(map[string]interface{})
			} else {
				dm = dv.(map[string]interface{})
			}
			copyMap(dm, v.(map[string]interface{}), config)
			dst[k] = dm
		case []string:
			if dv == nil {
				dst[k] = v
			} else {
				switch dv.(type) {
				case []string:
					sv := v.([]string)
					dv := dv.([]string)
					for _, svv := range sv {
						if !str.InSlice(svv, dv) {
							dv = append(dv, svv)
						}
					}
					dst[k] = dv
				}
			}
		case string:
			if config.IgnoreEmptyValue && v == "" {
				continue
			}
			dst[k] = v
		default:
			dst[k] = v
		}
	}
	return
}


func CopyMap(dst, src map[string]interface{}, opts ...func(*Config)) {
	config := &Config{}

	for _, opt := range opts {
		opt(config)
	}

	copyMap(dst, src, config)
}