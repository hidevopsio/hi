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

/*
Package inject implements dependency injection.

Dependency injection with the struct tag `inject:""` or the constructor.

Dependency injection in Go

Dependency injection is a concept valid for any programming language. The general concept behind dependency injection is
called Inversion of Control. According to this concept a struct should not configure its dependencies statically but
should be configured from the outside.

Dependency Injection design pattern allows us to remove the hard-coded dependencies and make our application loosely
coupled, extendable and maintainable.

Dependency Injection is the idea that your components (usually structs in go) should receive their dependencies when
being created.This runs counter to the associated anti-pattern of components building their own dependencies during
initialization.

A Go struct has a dependency on another struct, if it uses an instance of this struct. We call this a struct dependency.
For example, a struct which accesses a user controller has a dependency on user service struct.

Ideally Go struct should be as independent as possible from other Go struct. This increases the possibility of reusing
these struct and to be able to test them independently from other struct.

Dependency injection by constructor

To use dependency injection, first, you need to register the dependency in init func by calling
app.Component(newFooService), newFoo is the constructor of the dependency.

	// dependency foo
	type FooService struct {
	}

	// FooService constructor
	func newFooService() *FooService {
		return &Foo{}
	}

	func init() {
		app.Component(newFooService)
	}

	// the consumer barController that depends on FooService
	type barController {
		fooService *FooService
	}

	// the consumer's constructor newBarController that inject the instance of FooService
	func newBarController(fooService *FooService) *barController {
		return &barService{
			fooService: fooService,
		}
	}

Auto Configuration

Auto Configuration is another cool feature that comes out of the box with Hiboot,
for more details, please see https://godoc.org/github.com/hidevopsio/hiboot/pkg/starter

Example

The following example shows a struct which has no hard dependencies.

*/
package inject
