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

package cli_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/hidevopsio/hiboot/pkg/log"
	"sync"
	"testing"
	"time"
)

var mux = &sync.Mutex{}

func init() {
	log.SetLevel(log.DebugLevel)
}

type rootCommand struct {
	cli.RootCommand
	Profile string
	IntVal  int
	BoolVal bool
}

func newRootCommand(foo *fooCommand) *rootCommand {
	c := &rootCommand{}
	c.Use = "demo"
	c.Short = "demo command"
	c.Long = "Run demo command"

	pf := c.PersistentFlags()
	pf.StringVarP(&c.Profile, "profile", "p", "dev", "e.g. --profile=test")
	pf.IntVarP(&c.IntVal, "integer", "i", 0, "e.g. --integer=1")
	pf.BoolVarP(&c.BoolVal, "bool", "b", false, "e.g. --bool=true or -b")

	c.Add(foo)
	return c
}

func (c *rootCommand) Run(args []string) (err error) {
	log.Debugf("on demo command - profile: %v, intVal: %v, boolVal: %v", c.Profile, c.IntVal, c.BoolVal)
	return
}

type fooCommand struct {
	cli.SubCommand
}

func newFooCommand(bar *barCommand, baz *bazCommand) *fooCommand {
	c := new(fooCommand)
	c.Use = "foo"
	c.Short = "foo command"
	c.Long = "Run foo command"
	c.Add(bar, baz)
	return c
}

func (c *fooCommand) Run(args []string) (err error) {
	log.Debug("on foo command")
	return nil
}

func (c *fooCommand) OnDaz(args []string) bool {
	log.Debug("on daz command")
	return false
}

func (c *fooCommand) OnBuzz(args []string) bool {
	log.Debug("on buzz command")
	return true
}

func (c *fooCommand) OnFoo(args []string) error {
	log.Debug("on foo command")
	return nil
}

func (c *fooCommand) OnFooBar(args []string) error {
	log.Debug("on foo command")
	// just for the sake of testing
	return errors.New("testing on fooBar command")
}

type barCommand struct {
	cli.SubCommand
}

func newBarCommand() *barCommand {
	c := new(barCommand)
	c.Use = "bar"
	c.Short = "bar command"
	c.Long = "Run bar command"

	return c
}

func (c *barCommand) Run(args []string) (err error) {
	log.Debug("on bar command")
	return nil
}

type bazCommand struct {
	cli.SubCommand
}

func newBazCommand() *bazCommand {
	c := new(bazCommand)
	c.Use = "baz"
	c.Short = "baz command"
	c.Long = "Run baz command"
	return c
}

func (c *bazCommand) Run(args []string) (err error) {
	log.Debug("on baz command")
	return nil
}

// demo foo bar
func TestCliApplication(t *testing.T) {
	mux.Lock()

	app.Register(newFooCommand, newBarCommand, newBazCommand)
	testApp := cli.NewTestApplication(t, newRootCommand).
		SetProperty("foo", "bar")

	t.Run("should run root command", func(t *testing.T) {
		_, err := testApp.Run("-p", "test", "-i", "2")
		assert.Equal(t, nil, err)
	})

	t.Run("should run foo command", func(t *testing.T) {
		_, err := testApp.Run("foo")
		assert.Equal(t, nil, err)
	})

	t.Run("should run bar command", func(t *testing.T) {
		_, err := testApp.Run("foo", "bar")
		assert.Equal(t, nil, err)
	})

	t.Run("should run baz command", func(t *testing.T) {
		_, err := testApp.Run("foo", "baz")
		assert.Equal(t, nil, err)
	})

	t.Run("should run baz command", func(t *testing.T) {
		_, err := testApp.Run("foo", "daz")
		assert.Equal(t, nil, err)
	})

	t.Run("should run baz command", func(t *testing.T) {
		_, err := testApp.Run("foo", "buzz")
		assert.Equal(t, nil, err)
	})

	t.Run("should run foo foo command", func(t *testing.T) {
		out, err := testApp.Run("foo", "foo")
		log.Debugf("%v", out)
		assert.Equal(t, nil, err)
	})

	mux.Unlock()
}

// demo foo bar
func TestCliFoo(t *testing.T) {
	mux.Lock()
	app.Register(newFooCommand, newBarCommand, newBazCommand)
	testApp := cli.NewTestApplication(t, newRootCommand).
		SetProperty("foo", "bar")

	t.Run("should run foo foo command", func(t *testing.T) {
		out, err := testApp.Run("foo", "fooBar")
		log.Debugf("%v %v", out, err)
		//assert.NotEqual(t, nil, err)
	})
	mux.Unlock()
}

func TestNewApplication(t *testing.T) {
	go cli.NewApplication().
		SetProperty("app.project", "cli-test-app").
		SetAddCommandLineProperties(false).
		Run()
	time.Sleep(2 * time.Second)
}

type A struct {
	Name string
}

func (a *A) Run(x string, y int) {
	log.Debugf("name: %v, x: %v, y: %v", a.Name, x, y)
}

type B struct {
	A
}

func (b *B) Run(x string) {
	log.Debugf("x: %v", x)
	b.A.Run(x, 123)
}

func TestAB(t *testing.T) {
	b := new(B)
	b.Name = "bb"
	b.Run("hello")
}



type circularFoo struct {
	circularBar *circularBar
}

func newCircularFoo(circularBar *circularBar) *circularFoo {
	return &circularFoo{
		circularBar: circularBar,
	}
}

type circularBar struct {
	circularFoo *circularFoo
}

func newCircularBar(circularFoo *circularFoo) *circularBar {
	return &circularBar{
		circularFoo: circularFoo,
	}
}

type circularDiCommand struct {
	cli.RootCommand

	circularFoo *circularFoo
}

func newCircularCommand(circularFoo *circularFoo) *circularDiCommand {
	c := new(circularDiCommand)
	c.Use = "circular"
	c.Short = "circular command"
	c.Long = "Run circular command"

	c.circularFoo = circularFoo
	return c
}

func (c *circularDiCommand) Run(args []string) (err error) {
	log.Debug("on circular command")
	return nil
}

func TestApplicationWithCircularDI(t *testing.T) {
	testApp := cli.NewTestApplication(t, newCircularCommand)

	t.Run("should detect circular di", func(t *testing.T) {
		_, err := testApp.Run("circular")
		assert.Equal(t, nil, err)
	})
}

