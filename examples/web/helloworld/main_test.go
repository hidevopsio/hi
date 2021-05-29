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
	"time"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/swagger"
)

var mu sync.Mutex
func TestRunMain(t *testing.T) {
	mu.Lock()
	go main()
	mu.Unlock()
}

func TestController(t *testing.T) {
	mu.Lock()

	time.Sleep(2 * time.Second)

	app.Register(swagger.ApiInfoBuilder().
		Title("HiBoot Example - Hello world").
		Description("This is an example that demonstrate the basic usage"))

	web.NewTestApp(t, new(Controller)).
		SetProperty("server.port", "8081").
		SetProperty(app.ProfilesInclude, swagger.Profile, web.Profile, actuator.Profile).
		Run(t).
		Get("/").
		Expect().Status(http.StatusOK)

	mu.Unlock()
}
