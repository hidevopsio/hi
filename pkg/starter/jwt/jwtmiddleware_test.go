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

package jwt_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
	"testing"
	"time"
)

type FakeContext struct {
}

func (c *FakeContext) Next() {
	log.Debug("FakeContext.Next()")
}

func (c *FakeContext) StopExecution() {
	log.Debug("FakeContext.Next()")
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestCheckJWT(t *testing.T) {

	t.Run("should report error if jwt properties does not injected", func(t *testing.T) {

		jwtToken := jwt.NewJwtToken(&jwt.Properties{})

		assert.Equal(t, nil, jwtToken)
	})

	t.Run("should generate jwt token", func(t *testing.T) {

		jwtToken := jwt.NewJwtToken(&jwt.Properties{
			PrivateKeyPath: "config/ssl/app.rsa",
			PublicKeyPath:  "config/ssl/app.rsa.pub",
		})

		token, err := jwtToken.Generate(jwt.Map{
			"username": "johndoe",
			"password": "PA$$W0RD",
		}, 500, time.Millisecond)

		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, token)
	})

}
