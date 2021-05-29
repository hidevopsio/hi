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
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"io"
	"reflect"
	"strings"
)

// Command the command interface for cli application
type Command interface {
	// EmbeddedCommand return the embedded command
	EmbeddedCommand() *cobra.Command
	// Add add a new command
	Add(commands ...Command) Command
	// HasChild check if it has child
	HasChild() bool
	// Children get children
	Children() []Command
	// Exec execute the command
	Exec() error
	// GetName get the command name
	GetName() string
	// FullName get the command full name
	FullName() string
	// SetName set the command name
	SetName(name string) Command
	// SetFullName set the command full name
	SetFullName(name string) Command
	// Parent get parent command
	Parent() Command
	// SetParent set parent command
	SetParent(p Command) Command
	// Run the callback of command, once the command is received, this method will be called
	Run(args []string) error
	// Find find child command
	Find(name string) (Command, error)
	// SetOutput set output, type: io.Writer
	SetOutput(output io.Writer)
	// SetArgs set args
	SetArgs(a []string)
	// ExecuteC execute command
	ExecuteC() (cmd *cobra.Command, err error)
	// PersistentFlags get persistent flags
	PersistentFlags() *pflag.FlagSet
}

const (
	actionPrefix = "On"
)

// ErrCommandNotFound command not found error
var ErrCommandNotFound = errors.New("command not found")

// ErrCommandHandlerNotFound the error message for 'command handler not found'
var ErrCommandHandlerNotFound = errors.New("command handler not found")

// baseCommand is the base command
type baseCommand struct {
	cobra.Command
	name     string
	fullName string
	parent   Command
	children []Command
}

// RootCommand root command qualifier
type RootCommand struct {
	at.Qualifier `value:"cli.rootCommand"`
	baseCommand
}

// SubCommand sub command qualifier
type SubCommand struct {
	baseCommand
}

// Dispatch method with OnAction prefix
func Dispatch(c Command, args []string) (retVal interface{}, err error) {
	if len(args) > 0 && args[0] != "" {
		methodName := actionPrefix + strings.Title(args[0])
		retVal, err = reflector.CallMethodByName(c, methodName, args[1:])
		return
	}
	return nil, ErrCommandHandlerNotFound
}

// Register register
func Register(c Command) {
	var next bool
	c.EmbeddedCommand().RunE = func(cmd *cobra.Command, args []string) error {
		result, err := Dispatch(c, args)
		if err == nil && result != nil {
			typ := reflect.TypeOf(result)
			typName := typ.Name()
			switch typName {
			case types.Bool:
				next = result.(bool)
			default:
				// assume that error is the default return type
				return result.(error)
			}
		} else {
			next = true
		}

		if next {
			return c.Run(args)
		}
		return nil
	}
}

// EmbeddedCommand get embedded command
func (c *baseCommand) EmbeddedCommand() *cobra.Command {
	return &c.Command
}

// Run method
func (c *baseCommand) Run(args []string) error {
	return nil
}

// Exec exec method
func (c *baseCommand) Exec() error {

	// UnknownFlags will ignore unknown flags errors and continue parsing rest of the flags
	c.FParseErrWhitelist.UnknownFlags = true

	return c.Execute()
}

// HasChild check whether it has child or not
func (c *baseCommand) HasChild() bool {
	return len(c.children) > 0
}

// Children get children
func (c *baseCommand) Children() []Command {
	return c.children
}

func (c *baseCommand) addChild(child Command) {
	if child.GetName() == "" {
		name := reflector.ParseObjectName(child, "Command")
		child.SetName(name)
	}
	Register(child)
	child.SetParent(c)
	c.children = append(c.children, child)
	c.AddCommand(child.EmbeddedCommand())
}

// Add added child command
func (c *baseCommand) Add(commands ...Command) Command {
	for _, command := range commands {
		c.addChild(command)
	}
	return c
}

// GetName get command name
func (c *baseCommand) GetName() string {
	return c.name
}

// SetName set command name
func (c *baseCommand) SetName(name string) Command {
	c.name = name
	return c
}

// FullName get command full name
func (c *baseCommand) FullName() string {
	if c.fullName == "" {
		c.fullName = c.name
	}
	return c.fullName
}

// SetFullName set command full name
func (c *baseCommand) SetFullName(name string) Command {
	c.fullName = name
	return c
}

// Parent get parent command
func (c *baseCommand) Parent() Command {
	return c.parent
}

// SetParent set parent command
func (c *baseCommand) SetParent(p Command) Command {
	c.parent = p
	return c
}

// Find find child command
func (c *baseCommand) Find(name string) (Command, error) {
	if c.name == name {
		return c, nil
	}

	for _, cmd := range c.children {
		if name == cmd.GetName() {
			return cmd, nil
		}
	}
	return nil, ErrCommandNotFound
}
