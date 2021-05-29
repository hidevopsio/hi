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

// Package logging provides the hiboot starter for injectable logging dependency
package logging

import (
	"github.com/kataras/iris/middleware/logger"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
)

const (
	// Profile is the profile of logging, it should be as same as the package name
	Profile = "logging"
	// Level is the property for setting logging level
	Level = "logging.level"
	// LevelDebug is the logging level options
	LevelDebug = "debug"
	// LevelInfo is the logging level options
	LevelInfo = "info"
	// LevelWarn is the logging level options
	LevelWarn = "warn"
	// LevelError is the logging level options
	LevelError = "error"
	// LevelFatal is the logging level options
	LevelFatal = "fatal"
	// LevelDebug is the logging level options
	LevelDisable = "disable"
)

type configuration struct {
	app.Configuration
	Properties         *properties
	applicationContext app.ApplicationContext
}

func newConfiguration(applicationContext app.ApplicationContext) *configuration {
	return &configuration{
		applicationContext: applicationContext,
	}
}

func init() {
	app.Register(newConfiguration)
}

// LoggerHandler config logger handler
func (c *configuration) LoggerHandler() context.Handler {
	loggerHandler := logger.New(logger.Config{
		// Status displays status code
		Status: c.Properties.Status,
		// IP displays request's remote address
		IP: c.Properties.IP,
		// Method displays the http method
		Method: c.Properties.Method,
		// Path displays the request path
		Path: c.Properties.Path,
		// Query appends the url query to the Path.
		Query: c.Properties.Query,

		Columns: c.Properties.Columns,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: c.Properties.ContextKeys,

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: c.Properties.HeaderKeys,
	})

	lh := context.NewHandler(loggerHandler)

	c.applicationContext.Use(lh)

	return lh
}
