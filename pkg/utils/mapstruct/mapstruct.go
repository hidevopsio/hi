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

//Package mapstruct provides utilities that decode map and inject values into struct
package mapstruct

import (
	"encoding/json"
	"fmt"
	"github.com/hidevopsio/mapstructure"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"reflect"
)


func WithAnnotation(config *mapstructure.DecoderConfig) {
	config.TagName = "at"
}

func WithSquash(config *mapstructure.DecoderConfig) {
	config.Squash = true
}

func WithWeaklyTypedInput(config *mapstructure.DecoderConfig) {
	config.WeaklyTypedInput = true
}

// Decode decode (convert) map to struct
func Decode(to interface{}, from interface{}, opts ...func (*mapstructure.DecoderConfig) ) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           to,
		TagName:          "json",
	}

	for _, opt := range opts {
		opt(config)
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if from == nil {
		return fmt.Errorf("parameters of mapstruct.Decode must not be nil")
	}

	return decoder.Decode(from)
}

// DecodeStructToMap
// TODO: should improve the performance
func DecodeStructToMap(val interface{}) (sm map[string]interface{}, ok bool) {
	sv := reflector.IndirectValue(val)
	sk := sv.Kind()
	if sk == reflect.Struct {
		// convert object to []byte
		bs, err := json.Marshal(val)
		if err == nil {
			// define new map sm, unmarshal bs to sm
			err = json.Unmarshal(bs, &sm)
			ok = true
		}
	}
	return
}