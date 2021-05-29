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

package inject_test

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
)

//This example shows that the dependency is injected through the constructor
func Example() {
	web.NewApplication().Run()
}

// HelloService is a simple service interface, with interface, we can mock a fake service in unit test
type HelloService interface {
	SayHello(name string) string
}

type helloServiceImpl struct {
}

func init() {
	// Register Rest Controller through constructor newHelloController
	// Register Service through constructor newHelloService
	app.Register(newHelloController, newHelloService)
}

// please note that the return type name of the constructor HelloService,
// hiboot will instantiate a instance named helloService for dependency injection
func newHelloService() HelloService {
	return &helloServiceImpl{}
}

// SayHello is a service method implementation
func (s *helloServiceImpl) SayHello(name string) string {
	return "Hello" + name
}

// PATH: /login
type helloController struct {
	at.RestController
	helloService HelloService
}

// newHelloController inject helloService through the argument helloService HelloService on constructor
func newHelloController(helloService HelloService) *helloController {
	return &helloController{
		helloService: helloService,
	}
}

// Get /
// The first word of method name is the http method GET
func (c *helloController) Get(name string) string {
	return c.helloService.SayHello(name)
}
