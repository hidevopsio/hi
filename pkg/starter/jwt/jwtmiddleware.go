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
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	mwjwt "github.com/iris-contrib/middleware/jwt"
	ictx "github.com/kataras/iris/context"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/log"
)

// Middleware derived from github.com/iris-contrib/middleware/jwt/Middleware
type Middleware struct {
	mwjwt.Middleware
}

// Serve the middleware's action
func (m *Middleware) Serve(ctx context.Context) {
	log.Debug("Serve()")
	if err := m.CheckJWT(ctx); err != nil {
		ctx.ResponseError(err.Error(), http.StatusUnauthorized)
		ctx.StopExecution()
		return
	}
	// If everything ok then call next.
	ctx.Next()
}

// CheckJWT the main functionality, checks for token
func (m *Middleware) CheckJWT(ctx ictx.Context) error {
	log.Debug("CheckJWT()")
	if !m.Config.EnableAuthOnOptions {
		if ctx.Method() == http.MethodOptions {
			return nil
		}
	}

	// Use the specified token extractor to extract a token from the request
	token, err := m.Config.Extractor(ctx)

	if token == "" {
		// no token, try extract token form param
		token = ctx.URLParam("token")
		// TODO: Make this a bit more robust, parsing-wise
		authHeaderParts := strings.Split(token, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			log.Errorf("error extracting JWT: %v", token)
		} else {
			token = authHeaderParts[1]
		}
	}

	// If debugging is turned on, log the outcome
	if err != nil {
		log.Errorf("error extracting JWT: %v", err)
	} else {
		log.Debugf("token extracted: %s", token)
	}

	// If an error occurs, call the error handler and return an error
	if err != nil {
		return fmt.Errorf("error extracting token: %v", err)
	}

	// If the token is empty...
	if token == "" {
		// Check if it was required
		if m.Config.CredentialsOptional {
			log.Debug("  No credentials found (CredentialsOptional=true)")
			// No error, just no token (and that is ok given that CredentialsOptional is true)
			return nil
		}

		// If we get here, the required token is missing
		errorMsg := "Required authorization token not found"
		log.Debug("  Error: No credentials found (CredentialsOptional=false)")
		return fmt.Errorf(errorMsg)
	}

	// Now parse the token
	parsedToken, err := jwt.Parse(token, m.Config.ValidationKeyGetter)
	// Check if there was an error in parsing...
	if err != nil || !parsedToken.Valid {
		log.Debugf("Error parsing token: %v", err)
		return fmt.Errorf("error parsing token: %v", err)
	}

	if m.Config.SigningMethod != nil && m.Config.SigningMethod.Alg() != parsedToken.Header["alg"] {
		message := fmt.Sprintf("Expected %s signing method but token specified %s",
			m.Config.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		log.Debugf("Error validating token algorithm: %s", message)
		return fmt.Errorf("error validating token algorithm: %s", message)
	}

	log.Debugf("JWT: %v", parsedToken)

	// If we get here, everything worked and we can set the
	// user property in context.
	ctx.Values().Set(m.Config.ContextKey, parsedToken)

	return nil
}

// NewJwtMiddleware New constructs a new Secure instance with supplied options.
func NewJwtMiddleware(cfg ...mwjwt.Config) *Middleware {

	var c mwjwt.Config
	if len(cfg) != 0 {
		c = cfg[0]
	}

	if c.ContextKey == "" {
		c.ContextKey = mwjwt.DefaultContextKey
	}

	if c.ErrorHandler == nil {
		c.ErrorHandler = mwjwt.OnError
	}

	if c.Extractor == nil {
		c.Extractor = mwjwt.FromAuthHeader
	}

	return &Middleware{mwjwt.Middleware{Config: c}}
}
