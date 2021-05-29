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

// Package fake provides fake.ApplicationContext for unit testing
package fake

import "github.com/hidevopsio/hiboot/pkg/app/web/context"

// ApplicationContext application context
type ApplicationContext struct {
}

// RegisterController register controller by interface
func (a *ApplicationContext) RegisterController(controller interface{}) error {
	return nil
}

// Use use middleware
func (a *ApplicationContext) Use(handlers ...context.Handler) {

}

// GetProperty get application property by name
func (a *ApplicationContext) GetProperty(name string) (value interface{}, ok bool) {
	return
}

// GetInstance get application instance by name
func (a *ApplicationContext) GetInstance(params ...interface{}) (instance interface{}) {
	return
}
