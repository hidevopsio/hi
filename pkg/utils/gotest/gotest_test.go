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

package gotest

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestIsRunning(t *testing.T) {
	isTestRunning := IsRunning()

	assert.Equal(t, true, isTestRunning)
}

func TestParseArgs(t *testing.T) {
	args := []string{"foo", "bar", "baz"}
	log.Debug(flag.Args())
	ParseArgs(args)
	log.Debug(flag.Args())
}
