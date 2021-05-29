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

package jwt

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
)

type postProcessor struct {
	applicationContext app.ApplicationContext
	jwtMiddleware      *Middleware
}

func init() {
	// register postProcessor
	app.RegisterPostProcessor(newPostProcessor)
}

func newPostProcessor(applicationContext app.ApplicationContext, jwtMiddleware *Middleware) *postProcessor {
	return &postProcessor{
		applicationContext: applicationContext,
		jwtMiddleware:      jwtMiddleware,
	}
}

// AfterInitialization jwt post processing
func (p *postProcessor) AfterInitialization() {
	//log.Debug("[jwt] AfterInitialization")

	// use jwt
	p.applicationContext.Use(p.jwtMiddleware.Serve)

	// finally register jwt controllers
	_ = p.applicationContext.RegisterController(at.JwtRestController{})
}
