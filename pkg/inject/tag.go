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
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"reflect"
	"strings"
)

// Tag the interface of Tag
type Tag interface {
	// Init init tag
	Init(configurableFactory factory.InstantiateFactory)
	// Decode parse tag and do dependency injection
	Decode(object reflect.Value, field reflect.StructField, property string) (retVal interface{})
	// Properties get properties
	Properties() cmap.ConcurrentMap
	// IsSingleton check if it is Singleton
	IsSingleton() bool
}

// BaseTag is the base struct of tag
type BaseTag struct {
	instantiateFactory factory.InstantiateFactory
	properties         cmap.ConcurrentMap
	systemConfig       *system.Configuration
}

// IsSingleton check if it is Singleton
func (t *BaseTag) IsSingleton() bool {
	return false
}

// Init init the tag
func (t *BaseTag) Init(configurableFactory factory.InstantiateFactory) {
	t.instantiateFactory = configurableFactory
}

// ParseProperties parse properties
func (t *BaseTag) ParseProperties(tag string) cmap.ConcurrentMap {
	t.properties = cmap.New()

	args := strings.Split(tag, ",")
	for _, v := range args {
		//log.Debug(v)
		n := strings.Index(v, "=")
		if n > 0 {
			key := v[:n]
			val := v[n+1:]
			if key != "" && val != "" {
				// check if val contains reference or env
				// TODO: should lookup certain config instead of for loop
				replacedVal := t.instantiateFactory.Replace(val)
				t.properties.Set(key, replacedVal)
			}
		}
	}
	return t.properties
}

// Properties get properties
func (t *BaseTag) Properties() cmap.ConcurrentMap {
	return t.properties
}

// Decode no implementation for base tag
func (t *BaseTag) Decode(object reflect.Value, field reflect.StructField, property string) (retVal interface{}) {
	return nil
}

