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
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
	"net/http"
	"strings"
	"testing"
	"time"
)

type userRequest struct {
	model.RequestBody
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type fooController struct {
	at.RestController
	token    jwt.Token
	tokenStr string
}

// for test only: token will expired in 1 second
var tokenExpiredSecond = int64(1)

func init() {
	log.SetLevel(log.DebugLevel)
}

func newFooController(token jwt.Token) *fooController {
	return &fooController{
		token: token,
	}
}

func (c *fooController) Get() string {
	log.Debug("fooController.Get")

	return "Hello, world"
}

func (c *fooController) PostLogin(request *userRequest) (response model.Response, err error) {
	log.Debug("fooController.Login")

	// you make validate username and password first
	token, err := c.token.Generate(jwt.Map{
		"username": request.Username,
		"password": request.Password,
	}, tokenExpiredSecond, time.Second)
	response = new(model.BaseResponse)
	response.SetData(token)
	c.tokenStr = token

	return
}

// barController
type barController struct {
	at.JwtRestController
}

func newBarController() *barController {
	return &barController{}
}

func (c *barController) Before(ctx context.Context, tokenProperties *jwt.TokenProperties) {
	log.Debug("barController.Before")

	_, ok := tokenProperties.GetAll()
	if !ok {
		return
	}

	_, ok = tokenProperties.Items()
	if !ok {
		return
	}

	username := tokenProperties.Get("username")
	if username == "" {
		return
	}

	ctx.Next()
}

func (c *barController) Get() string {
	log.Debug("barController.Get")

	return "Hello, world"
}

func (c *barController) Options() {
	log.Debug("barController.Options")
}

func TestJwtController(t *testing.T) {
	testApp := web.NewTestApp(newFooController, newBarController).SetProperty(app.ProfilesInclude, web.Profile, jwt.Profile).Run(t)
	appContext := testApp.(app.ApplicationContext)

	fc := appContext.GetInstance(fooController{})
	assert.NotEqual(t, nil, fc)
	fooCtrl := fc.(*fooController)

	bc := appContext.GetInstance(barController{})
	assert.NotEqual(t, nil, bc)

	t.Run("should login with POST /foo/login", func(t *testing.T) {
		testApp.
			Post("/foo/login").
			WithJSON(&userRequest{Username: "johndoe", Password: "iHop91#15"}).
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusUnauthorized after GET /bar", func(t *testing.T) {
		testApp.Get("/bar").
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should return http.StatusOK after GET /bar", func(t *testing.T) {
		testApp.Options("/bar").
			Expect().Status(http.StatusOK)
	})

	token := fmt.Sprintf("Bearer %v", fooCtrl.tokenStr)
	t.Run("should return http.StatusOK on /bar with jwt token", func(t *testing.T) {

		testApp.Get("/bar").
			WithQuery("token", token).
			Expect().Status(http.StatusOK)

		testApp.Get("/bar").
			WithHeader("Authorization", token).
			Expect().Status(http.StatusOK)

		testApp.Get("/bar").
			WithHeader("Authorization", "1df%^!"+token).
			Expect().Status(http.StatusUnauthorized)

		testApp.Get("/bar").
			WithHeader("Authorization", "").
			Expect().Status(http.StatusUnauthorized)
	})

	applicationContext := testApp.(app.ApplicationContext)
	jwtMiddleware := applicationContext.GetInstance(jwt.Middleware{}).(*jwt.Middleware)

	t.Run("should report error with SigningMethod", func(t *testing.T) {
		jwtMiddleware.Config.SigningMethod = jwtgo.SigningMethodRS512
		testApp.Get("/bar").
			WithHeader("Authorization", token).
			Expect().Status(http.StatusUnauthorized)
		jwtMiddleware.Config.SigningMethod = jwtgo.SigningMethodRS256
	})

	t.Run("should failed with invalid token", func(t *testing.T) {
		testApp.Get("/bar").
			WithHeader("Authorization", strings.Replace(token, "1", "2", -1)).
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should set CredentialsOptional", func(t *testing.T) {
		time.Sleep(2 * time.Second)

		testApp.Get("/bar").
			WithHeader("Authorization", token).
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should set CredentialsOptional", func(t *testing.T) {
		jwtMiddleware.Config.CredentialsOptional = true
		testApp.Get("/bar").
			WithHeader("Authorization", "").
			Expect().Status(http.StatusOK)
	})

}

