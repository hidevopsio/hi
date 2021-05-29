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

// Package jaeger provides the quick start jaeger application example
// main package
package main

// import web starter from hiboot
import (
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/jaeger"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"time"
)

// Controller Rest Controller with path /
// RESTful Controller, derived from at.RestController. The context mapping of this controller is '/' by default
type Controller struct {
	// at.RestController or at.RestController must be embedded here
	at.RestController
}

// Get GET /formatter/{format}
func (c *Controller) GetByFormatter(formatter string, span *jaeger.ChildSpan) string {
	defer span.Finish()
	greeting := span.BaggageItem("greeting")
	if greeting == "" {
		greeting = "Hello"
	}

	helloStr := fmt.Sprintf("[%s] %s, %s", time.Now().Format(time.Stamp), greeting, formatter)

	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	// response
	return helloStr
}

func newController() *Controller {
	return &Controller{}
}

func init() {
	app.Register(newController)
}

// main function
func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.ProfilesInclude, logging.Profile, web.Profile, jaeger.Profile).
		Run()
}
