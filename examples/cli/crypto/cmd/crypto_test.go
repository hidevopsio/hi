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
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/hidevopsio/hiboot/pkg/log"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestCryptoCommands(t *testing.T) {
	testApp := cli.NewTestApplication(t, newCryptoCommand)

	t.Run("should run crypto rsa -e", func(t *testing.T) {
		_, err := testApp.Run("rsa", "-e", "-s", "hello")
		assert.Equal(t, nil, err)
	})

	t.Run("should run crypto rsa -d ", func(t *testing.T) {
		_, err := testApp.Run("rsa", "-d", "-s", "Rprrfl5LX9NRmWKEqJW8ckObVjznnMmq8i7x6Pv6n1GSoEL9dUomNKOr6Pgj7RuVzCc/I7Hya20BZO1PbzTquBMp/G5rcF2Vy7HF1UKr8buHtppB+n3ycTxFvPxQB2vMvLyMtDBc29QtGe3HHD8TS+3h1pSK5WZS+CMKPHT4sho=")
		assert.Equal(t, nil, err)
	})
}
