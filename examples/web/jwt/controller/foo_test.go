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

// Line 1: main package
package controller

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"net/http"
	"testing"
)

func TestFooController(t *testing.T) {

	testApp := web.RunTestApplication(t, newFooController)

	t.Run("should pass GET /foo", func(t *testing.T) {
		testApp.
			Get("/foo").
			WithQueryObject(fooRequestParam{Name: "Peter", Age: 18}).
			Expect().Status(http.StatusOK).Body().Contains("jwt-example-v1")
	})

	t.Run("should pass POST /foo", func(t *testing.T) {
		testApp.
			Post("/foo").
			WithJSON(fooRequestBody{Name: "Mike", Age: 18}).
			Expect().Status(http.StatusOK)
	})
}
