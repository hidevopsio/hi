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
	"crypto/rsa"
	"github.com/hidevopsio/hiboot/pkg/log"
	"time"

	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"io/ioutil"
)

// Map is the JWT map
type Map map[string]interface{}

// Token is the token interface
type Token interface {
	Generate(payload Map, expired int64, unit time.Duration) (string, error)
	VerifyKey() *rsa.PublicKey
}

type jwtToken struct {
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
	//jwtMiddleware *JwtMiddleware
	jwtEnabled bool
}

// NewJwtToken create new jwt token
func NewJwtToken(p *Properties) (token Token) {
	jt := new(jwtToken)
	err := jt.Initialize(p)
	if err == nil {
		token = jt
	}
	return
}

func (t *jwtToken) Initialize(p *Properties) (err error) {
	if io.IsPathNotExist(p.PrivateKeyPath) {
		err = fmt.Errorf("private key file %v does not exist", p.PrivateKeyPath)
	}
	if err == nil {
		var verifyBytes []byte
		signBytes, err := ioutil.ReadFile(p.PrivateKeyPath)
		if err == nil {
			t.signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
			if err == nil {
				verifyBytes, err = ioutil.ReadFile(p.PublicKeyPath)
				if err == nil {
					t.verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
					if err == nil {
						t.jwtEnabled = true
					}
				}
			}
		}
	}
	if err != nil {
		log.Error(err)
	}
	return err
}

func (t *jwtToken) VerifyKey() *rsa.PublicKey {
	return t.verifyKey
}

// Generate generates JWT token with specified exired time
func (t *jwtToken) Generate(payload Map, expired int64, unit time.Duration) (tokenString string, err error) {
	if t.jwtEnabled {
		claim := jwt.MapClaims{
			"exp": time.Now().Add(unit * time.Duration(expired)).Unix(),
			"iat": time.Now().Unix(),
		}

		for k, v := range payload {
			claim[k] = v
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err = token.SignedString(t.signKey)
	} else {
		err = fmt.Errorf("JWT is not initialized")
	}
	return
}
