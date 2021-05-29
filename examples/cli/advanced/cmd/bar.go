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

package cmd

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type barCommand struct {
	cli.SubCommand
}

func init() {
	app.Register(newBarCommand)
}

func newBarCommand() *barCommand {
	c := new(barCommand)
	c.Use = "bar"
	c.Short = "bar command"
	c.Long = "Run bar command"
	return c
}

// OnBaz run command bar baz, return true means it won't run next action, in this case is method Run(args []string)
func (c *barCommand) OnBaz(args []string) bool {
	log.Infof("on baz command")
	return true
}

// OnBuz run command bar buz, return true means it won't run next action, in this case is method Run(args []string)
func (c *barCommand) OnBuz(args []string) bool {
	log.Infof("on buz command")
	return true
}

// Run run bar command
func (c *barCommand) Run(args []string) error {
	log.Info("handle bar command")
	return nil
}
