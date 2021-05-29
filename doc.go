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

/*Package hiboot is a web/cli app application framework

Hiboot is a cloud native web and cli application framework written in Go.

Hiboot integrates the popular libraries but make them simpler, easier to use.
It borrowed some of the Spring features like dependency injection, aspect oriented programming, and auto configuration.
You can integrate any other libraries easily by auto configuration with dependency injection support. hiboot-data is the
typical project that implement customized hiboot starters. see https://godoc.org/github.com/hidevopsio/hiboot-data

Overview

	Web MVC - (Model-View-Controller).
	Auto Configuration - pre-create instance with properties configs for dependency injection.
	Dependency injection - with struct tag name `inject:""` or Constructor func.

Introduction to Hiboot

One of the most significant feature of Hiboot is Dependency Injection. Hiboot implements JSR-330 standard.

Let's say that we have two implementations of AuthenticationService, below will explain how does Hiboot work.

	type AuthenticationService interface {
		Authenticate(credential Credential) error
	}

	type basicAuthenticationService struct {
	}

	func newBasicAuthenticationService() AuthenticationService {
		return &basicAuthenticationService{}
	}

	func (s *basicAuthenticationService) Authenticate(credential Credential) error {
		// business logic ...
		return nil
	}

	type oauth2AuthenticationService struct {
	}

	func newOauth2AuthenticationService() AuthenticationService {
		return &oauth2AuthenticationService{}
	}

	func (s *oauth2AuthenticationService) Authenticate(credential Credential) error {
		// business logic ...
		return nil
	}

	func init() {
		app.Register(newBasicAuthenticationService, newOauth2AuthenticationService)
	}

Field Injection

In Hiboot the injection into fields is triggered by `inject:""` struct tag. when inject tag is present
on a field, Hiboot tries to resolve the object to inject by the type of the field. If several implementations
of the same service interface are available, you have to disambiguate which implementation you want to be
injected. This can be done by naming the field to specific implementation.

	type userController struct {
		at.RestController

		BasicAuthenticationService AuthenticationService	`inject:""`
		Oauth2AuthenticationService AuthenticationService	`inject:""`
	}

	func newUserController() {
		return &userController{}
	}

	func init() {
		app.Register(newUserController)
	}

Constructor Injection

Although Field Injection is pretty convenient, but the Constructor Injection is the first-class citizen, we
usually advise people to use constructor injection as it has below advantages,

	* It's testable, easy to implement unit test.
	* Syntax validation, with syntax validation on most of the IDEs to avoid typo.
	* No need to use a dedicated mechanism to ensure required properties are set.

	type userController struct {
		at.RestController

		basicAuthenticationService AuthenticationService
	}

	// Hiboot will inject the implementation of AuthenticationService
	func newUserController(basicAuthenticationService AuthenticationService) {
		return &userController{
			basicAuthenticationService: basicAuthenticationService,
		}
	}

	func init() {
		app.Register(newUserController)
	}


Features
	App
		cli - command line application
		web - web application

	Starters
		actuator - health check
		locale - locale starter
		logging - customized logging settings
		jwt - jwt starter
		grpc - grpc application starter

	Tags
		inject - inject generic instance into object
		default - inject default value into struct object
		value - inject string value or references / variables into struct string field

	Utils
		cmap - concurrent map
		copier - copy between struct
		crypto - aes, base64, md5, and rsa encryption / decryption
		gotest - go test util
		idgen - twitter snowflake id generator
		io - file io util
		mapstruct - convert map to struct
		replacer - replacing stuct field value with references or environment variables
		sort - sort slice elements
		str - string util enhancement util
		validator - struct field validation


Getting started

This section will show you how to create and run a simplest hiboot application. Let’s get started!

Getting started with Hiboot web application

Get the source code

	go get -u github.com/hidevopsio/hiboot
	cd $GOPATH/src/github.com/hidevopsio/hiboot/examples/web/helloworld/

Source Code

*/
package hiboot
