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

// Package helloworld provides the quick start web application example
// main package
package main

// import web starter from hiboot
import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/swagger"
)

// Controller Rest Controller with path /
// RESTful Controller, derived from at.RestController. The context mapping of this controller is '/' by default
type Controller struct {
	// at.RestController or at.RestController must be embedded here
	at.RestController
	at.RequestMapping `value:"/"`
}

// Get GET /
func (c *Controller) Get(at struct {
	at.GetMapping `value:"/"`
	at.Operation  `id:"helloWorld" description:"This is hello world API"`
	at.Produces   `values:"text/plain"`
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"response status OK"`
			at.Schema `type:"string" description:"returns hello world message"`
		}
	}
}) string {
	// response
	return "Hello world"
}

// main function
func main() {
	app.Register(swagger.ApiInfoBuilder().
		Title("HiBoot Example - Hello world").
		Description("This is an example that demonstrate the basic usage"))

	// create new web application and run it
	web.NewApplication(new(Controller)).
		SetProperty(app.ProfilesInclude, swagger.Profile, web.Profile, actuator.Profile).
		Run()
}
