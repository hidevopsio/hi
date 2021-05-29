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

package model_test


import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"testing"
)

func TestResponse(t *testing.T) {
	var response model.Response

	response = new(model.BaseResponse)

	response.SetCode(200)
	assert.Equal(t, 200, response.GetCode())

	response.SetMessage("Hello world")
	assert.Equal(t, "Hello world", response.GetMessage())

	response.SetData("example")
	assert.Equal(t, "example", response.GetData())

	t.Run("should check implements specific interface", func(t *testing.T) {
		respInfo := new(model.BaseResponseInfo)
		ok := reflector.Implements(respInfo, new(model.ResponseInfo))
		assert.True(t, ok)
	})

	t.Run("should check implements specific interface", func(t *testing.T) {
		resp := new(model.BaseResponse)
		ok := reflector.Implements(resp, new(model.ResponseInfo))
		assert.True(t, ok)
	})

	t.Run("should check implements specific interface", func(t *testing.T) {
		resp := new(model.BaseResponseInfo)
		ok := reflector.Implements(resp, new(model.Response))
		assert.False(t, ok)
	})
}
