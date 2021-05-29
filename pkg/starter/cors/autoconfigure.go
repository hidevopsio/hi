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

// package cors provides the hiboot starter for injectable jwt dependency
package cors

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
)

const (
	// Profile is the profile of jwt, it should be as same as the package name
	Profile = "cors"
)

type configuration struct {
	at.AutoConfiguration

	Properties *Properties
}

func init() {
	app.Register(newConfiguration)
}

func newConfiguration() *configuration {
	return &configuration{}
}

type Middleware struct {
	context.Handler
}

func (c *configuration) Middleware(applicationContext app.ApplicationContext) (mw *Middleware) {
	mw = new(Middleware)
	mw.Handler = NewMiddleware(c.Properties)
	applicationContext.Use(mw.Handler)
	return
}
