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
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
)

type valueTag struct {
	at.Tag `value:"value"`
	BaseTag
}

func init() {
	AddTag(new(valueTag))
}

func (t *valueTag) Decode(object reflect.Value, field reflect.StructField, property string) (retVal interface{}) {
	tag, ok := field.Tag.Lookup(t.Tag.AtValue)
	if ok {
		//log.Debug(valueTag)

		// check if filed type is slice
		kind := field.Type.Kind()
		needConvert := true
		retVal = t.instantiateFactory.Replace(tag)
		switch kind {
		case reflect.Slice:
			if retVal != tag {
				needConvert = false
			}
		case reflect.String:
			needConvert = false
		}

		if needConvert {
			in := fmt.Sprintf("%v", retVal)
			retVal = str.Convert(in, kind)
		}
	}
	return retVal
}
