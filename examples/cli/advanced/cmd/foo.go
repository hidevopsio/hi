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
	"github.com/hidevopsio/hiboot/examples/cli/advanced/model"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type fooCommand struct {
	cli.SubCommand

	fooBar *model.Foo
}

func newFooCommand(fooBar *model.Foo) *fooCommand {
	c := &fooCommand{
		fooBar: fooBar,
	}
	c.Use = "foo"
	c.Short = "foo command"
	c.Long = "Run foo command"

	return c
}

func init() {
	app.Register(newFooCommand)
}

func (c *fooCommand) Run(args []string) error {
	log.Info("handle foo command")
	return nil
}
