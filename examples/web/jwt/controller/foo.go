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

package controller

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
)

type fooRequestBody struct {
	at.RequestBody

	Name    string `json:"name" validate:"required"`
	AppName string `json:"appName" default:"${app.name}"`
	Age     int    `json:"age"`
}

type fooRequestParam struct {
	at.RequestParams

	Name string `json:"name" validate:"required"`
	Age  int    `json:"age"`
}

type fooResponse struct {
	at.ResponseBody
	AppNameVersion string `json:"appNameVersion"`
	Greeting string `json:"greeting"`
	Age      int    `json:"age"`
}

type fooController struct {
	at.RestController

	AppNameVersion string `value:"${app.name}-${app.version}"`
}

// init - add &FooController{} to web application
func init() {
	app.Register(newFooController)
}

func newFooController() *fooController {
	return &fooController{}
}

// Before intercept all requests that coming into this controller
func (c *fooController) Before(ctx context.Context) {
	log.Debug("FooController.Before")
	ctx.Next()
}

// Post endpoint POST /foo
func (c *fooController) Post(request *fooRequestBody) (response model.Response, err error) {
	log.Debug("FooController.Post")
	log.Debug(request)
	response = new(model.BaseResponse)
	response.SetData(&fooResponse{
		Greeting: "Hello, " + request.Name,
		Age:      request.Age})

	return
}

// Get endpoint GET /foo
func (c *fooController) Get(request *fooRequestParam) (response model.Response, err error) {
	log.Debug("FooController.Get")
	response = new(model.BaseResponse)
	response.SetData(&fooResponse{
		AppNameVersion: c.AppNameVersion,
		Greeting: "Hello, " + request.Name,
		Age:      request.Age})

	return
}

// After interceptor of the controller
func (c *fooController) After(ctx context.Context) {
	log.Debug("FooController.After")
}
