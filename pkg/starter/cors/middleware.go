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

package cors

import (
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/iris-contrib/middleware/cors"
)

// NewMiddleware
func NewMiddleware(properties *Properties) (crs context.Handler) {

	options := cors.Options{
		AllowedOrigins: properties.AllowedOrigins,
		AllowedHeaders: properties.AllowedHeaders,
		AllowedMethods: properties.AllowedMethods,
		ExposedHeaders: properties.ExposedHeaders,
		AllowCredentials: properties.AllowCredentials,
		Debug: properties.Debug,
		OptionsPassthrough: properties.OptionsPassthrough,
		MaxAge: properties.MaxAge,
	}

	crs = context.NewHandler(cors.New(options))

	return
}
