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

package cli

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"os"
	"path/filepath"
	"strings"
)

// Application cli application interface
type Application interface {
	app.Application
}

type application struct {
	app.BaseApplication
	root Command
}

// CommandNameValue is the command name value pair
type CommandNameValue struct {
	Name    string
	Command interface{}
}

const (
	// RootCommandName the instance name of cli.rootCommand
	RootCommandName = "cli.rootCommand"
)

// NewApplication create new cli application
func NewApplication(cmd ...interface{}) Application {
	a := new(application)
	a.initialize(cmd...)
	return a
}

func (a *application) initialize(cmd ...interface{}) (err error) {
	if len(cmd) > 0 {
		app.Register(RootCommandName, cmd[0])
	}
	err = a.Initialize()
	return
}

// Init initialize cli application
func (a *application) build() (err error) {

	a.Build()

	a.PrintStartupMessages()

	basename := filepath.Base(os.Args[0])
	basename = strings.ToLower(basename)
	basename = strings.TrimSuffix(basename, ".exe")

	f := a.ConfigurableFactory()
	f.SetInstance(app.ApplicationContextName, a)

	// build auto configurations
	err = a.BuildConfigurations()
	//if err != nil {
	//	return
	//}

	// set root command
	r := f.GetInstance(RootCommandName)
	var root Command
	if r != nil {
		root = r.(Command)
		Register(root)
		a.root = root
		root.EmbeddedCommand().Use = basename
	}
	return
}

// SetProperty set application property
func (a *application) SetProperty(name string, value ...interface{}) app.Application {
	a.BaseApplication.SetProperty(name, value...)
	return a
}

// SetAddCommandLineProperties set add command line properties to be enabled or disabled
func (a *application) SetAddCommandLineProperties(enabled bool) app.Application {
	a.BaseApplication.SetAddCommandLineProperties(enabled)
	return a
}

// Initialize init application
func (a *application) Initialize() error {
	return a.BaseApplication.Initialize()
}

// Run run the cli application
func (a *application) Run() {
	a.build()
	//log.Debug(commandContainer)
	if a.root != nil {
		a.root.Exec()
	}
}
