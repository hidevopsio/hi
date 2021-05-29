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

package jwt

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"testing"
	"time"
)

func init() {
	log.SetLevel(log.DebugLevel)
	io.EnsureWorkDir(1, "config/application.yml")
}

func TestAutoConfigure(t *testing.T) {
	t.Run("should create new jwt middleware", func(t *testing.T) {
		config := &configuration{
			Properties: &Properties{
				PrivateKeyPath: "config/ssl/app.rsa",
				PublicKeyPath:  "config/ssl/app.rsa.pub",
			},
		}

		token := config.Token()
		assert.NotEqual(t, nil, token)
		mw := config.Middleware(token.(*jwtToken))
		assert.NotEqual(t, nil, mw)
	})

	t.Run("should report if jwt ssl does not exist", func(t *testing.T) {
		config := &configuration{
			Properties: &Properties{
				PrivateKeyPath: "does-not-exist",
				PublicKeyPath:  "does-not-exist",
			},
		}

		token := config.Token().(*jwtToken)
		assert.Equal(t, false, token.jwtEnabled)

		err := token.Initialize(config.Properties)
		assert.NotEqual(t, nil, err)

		_, err = token.Generate(Map{
			"username": "johndoe",
			"password": "PA$$W0RD",
		}, 10, time.Second)
		assert.NotEqual(t, nil, err)
	})

}
