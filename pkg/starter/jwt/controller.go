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
	"github.com/dgrijalva/jwt-go"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
)

// Controller is the base controller for jwt.RestController
type Controller struct {
	at.JwtRestController
}

// TokenProperties is the struct for parse jwt token properties
type TokenProperties struct {
	at.ContextAware
	context context.Context
}

// newTokenProperties is the constructor of TokenProperties
func newTokenProperties(context context.Context) *TokenProperties {
	return &TokenProperties{context: context}
}

// Parse is an util that parsing JWT token from jwt.MapClaims
func (p *TokenProperties) Parse(claims jwt.MapClaims, prop string) (retVal string) {
	val, ok := claims[prop]
	if ok {
		retVal = fmt.Sprintf("%v", val)
	}
	return
}

// Get is an util that parsing JWT token and return single property from jwt.MapClaims
func (p *TokenProperties) Get(propName string) (propVal string) {
	claims, ok := p.GetAll()
	if ok {
		propVal = p.Parse(claims, propName)
	}
	return
}

// GetAll is an util that parsing JWT token and return all properties from jwt.MapClaims
func (p *TokenProperties) GetAll() (propMap map[string]interface{}, ok bool) {
	var token *jwt.Token
	if p.context != nil {
		jwtVal := p.context.Values().Get("jwt")
		if jwtVal != nil {
			token = jwtVal.(*jwt.Token)
			if token.Claims != nil {
				var claims jwt.MapClaims
				if claims, ok = token.Claims.(jwt.MapClaims); ok && token.Valid {
					propMap = claims
				}
			}
		}
	}
	return
}

// Items is an util that parsing JWT token and return all properties in map from jwt.MapClaims
func (p *TokenProperties) Items() (propMap map[string]string, ok bool) {
	propMap = make(map[string]string)

	claims, ok := p.GetAll()
	if ok {
		for name, value := range claims {
			propMap[name] = fmt.Sprintf("%v", value)
		}
	}
	return
}
