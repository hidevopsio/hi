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
package controller

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"net/http"
	"testing"
	"time"
)

func genJwtToken(timeout int64) (token string) {
	jwtToken := jwt.NewJwtToken(&jwt.Properties{
		PrivateKeyPath: "config/ssl/app.rsa",
		PublicKeyPath:  "config/ssl/app.rsa.pub",
	})
	pt, err := jwtToken.Generate(jwt.Map{
		"username": "johndoe",
		"password": "PA$$W0RD",
	}, timeout, time.Second)
	if err == nil {
		token = fmt.Sprintf("Bearer %v", pt)
	}
	return
}

func TestBarWithToken(t *testing.T) {
	app := web.NewTestApp(t, newBarController).
		Run(t)

	log.SetLevel(log.DebugLevel)
	log.Println(io.GetWorkDir())
	token := genJwtToken(1)
	if token != "" {
		t.Run("should pass with jwt token", func(t *testing.T) {
			app.Get("/bar").
				WithHeader("Authorization", token).
				Expect().Status(http.StatusOK)
		})

		time.Sleep(2 * time.Second)

		t.Run("should not pass with expired jwt token", func(t *testing.T) {
			app.Get("/bar").
				WithHeader("Authorization", token).
				Expect().Status(http.StatusUnauthorized)
		})
	}
}
