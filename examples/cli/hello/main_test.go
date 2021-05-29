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

// Line 1: main package
package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"sync"
	"testing"
)

var mu sync.Mutex
func TestRunMain(t *testing.T) {
	mu.Lock()
	go main()
	mu.Unlock()
}

func TestHelloCommands(t *testing.T) {
	testApp := cli.NewTestApplication(t, newRootCommand)

	t.Run("should run hello command", func(t *testing.T) {
		_, err := testApp.Run("--to", "${app.name}-cmd")
		assert.Equal(t, nil, err)
	})
}
