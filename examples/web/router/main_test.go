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
package main

import (
	"net/http"
	"sync"
	"testing"

	"github.com/hidevopsio/hiboot/pkg/app/web"
)

var mu sync.Mutex
//func TestRunMain(t *testing.T) {
//	mu.Lock()
//	go main()
//	mu.Unlock()
//}

func TestController(t *testing.T) {
	mu.Lock()
	testApp := web.NewTestApp(t).SetProperty("server.context_path", "/router-example").Run(t)

	t.Run("should send post request to user", func(t *testing.T) {
		testApp.Post("/router-example/user").
			WithJSON(&UserRequest{UserVO: UserVO{Username: "John", Password: "password"}}).
			Expect().Status(http.StatusOK)
	})

	t.Run("should pass test for user controller", func(t *testing.T) {
		testApp.Delete("/router-example/user/123").
			Expect().Status(http.StatusOK)

		testApp.Get("/router-example/user").
			Expect().Status(http.StatusOK)

		testApp.Patch("/router-example/user/123").
			Expect().Status(http.StatusOK)
	})

	t.Run("should pass test for user controller with path variable", func(t *testing.T) {
		testApp.Get("/router-example/user/123/and/John").
			Expect().Status(http.StatusOK)
	})

	t.Run("should pass test for organization controller", func(t *testing.T) {
		testApp.Get("/router-example/organization/official-site").
			Expect().Status(http.StatusOK)
	})

	t.Run("should pass test for organization controller with path variable", func(t *testing.T) {
		testApp.Get("/router-example/organization/123/and/John").
			Expect().Status(http.StatusOK)

	})
	mu.Unlock()
}
