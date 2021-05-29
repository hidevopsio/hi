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
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/hidevopsio/hiboot/pkg/log"
)

// RootCommand is the root command
type RootCommand struct {
	cli.RootCommand

	profile string
	timeout int

	MagicNumber int `value:"${magic.number}"`
}

// NewRootCommand the root command
func NewRootCommand(second *secondCommand) *RootCommand {
	c := new(RootCommand)
	c.Use = "first"
	c.Short = "first command"
	c.Long = "Run first command"
	c.ValidArgs = []string{"baz"}
	pf := c.PersistentFlags()
	pf.StringVarP(&c.profile, "profile", "p", "dev", "e.g. --profile=test")
	pf.IntVarP(&c.timeout, "timeout", "t", 1, "e.g. --timeout=1")
	c.Add(second)
	return c
}

// Run root command handler
func (c *RootCommand) Run(args []string) error {
	log.Debugf("root: %v", args)
	log.Infof("handle first command: profile=%v, timeout=%v, magic.number=%v", c.profile, c.timeout, c.MagicNumber)
	return nil
}
