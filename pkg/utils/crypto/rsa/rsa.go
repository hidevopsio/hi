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

// Package rsa provides rsa encryption/decryption utilities
package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/hidevopsio/hiboot/pkg/utils/crypto"
	"github.com/hidevopsio/hiboot/pkg/utils/crypto/base64"
)

//openssl genrsa -out private.pem 1024
var defaultPrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDTyDSWdsUUZlym136QZnDgcg6FOmbak3Wkr85pxsoTzta9+ARo
xO8n/rw05ZtFsNGEj4ehFRt4+xU2v9GMFHtudj6GdtOs4rP6YpVbzGQd22K0tkAD
bsxMBHKsHsFcg02uFGSF6qIPnrWJIN0wNEJrNumo+XJ9EsXNggfvnp/t8QIDAQAB
AoGASOSG3ScCDFRzHWXchR0NSuNXBhok3qSUgFuWgyfN0+WEWhxsgBcQbcaqxtYk
jGcgpiy0tQfLzec11QNOv5UpKiqQU0TxYJjkRQW+Jz7J1K1YF4IokyJPVU3Rz/IB
RmtFMySZuP+uQdznoqPG5843pCcwCz5q5Gq89zgsubyMNjECQQDy9tex4dwrYtpU
E2abQ+mcOwaJlBMzvjN4sNqkd/ZpvtAfwwt5HE2qAYdNDFk6hI5i3HYQTtMD1+gq
vgWANegFAkEA3yUR7x8gNoSvULdDqOQmrsCqu59/HcIHGcnyZFXPUmMRAXq51aQE
8oS7dXFDUsUK1tdXxGDWYWX4E513kent/QJACIIaZZFKO363tJAFXNSQ/raWcQTt
czdq6AQRdAb7axKSiTo0UaZrFdP93/XZvhHcRpv/ymxoOU87QxvqZ2X73QJBAKbQ
jpylNy6qeGkt5729eZGQZNJIRP4ZC7fiuXr6jzd26cKiXYRxzmUChyUf3AVnWlgL
uggLoJhFY3Q+dqG1MH0CQFppAzwYv0NMmG9eFT2XNrMK7HDAAfPQU+hqd0M3Np7n
J6C0U/ErlLWE8GqXZP7+rPLGYacyUDJiZZMDB2X4AP0=
-----END RSA PRIVATE KEY-----
`)

//openssl rsa -in private.pem -pubout -out public.pem
var defaultPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDTyDSWdsUUZlym136QZnDgcg6F
Ombak3Wkr85pxsoTzta9+ARoxO8n/rw05ZtFsNGEj4ehFRt4+xU2v9GMFHtudj6G
dtOs4rP6YpVbzGQd22K0tkADbsxMBHKsHsFcg02uFGSF6qIPnrWJIN0wNEJrNumo
+XJ9EsXNggfvnp/t8QIDAQAB
-----END PUBLIC KEY-----
`)

// Encrypt to rsa string
func Encrypt(input []byte, publicKey ...[]byte) ([]byte, error) {
	//decrypt pem public key
	// TODO: check if the publicKey or privateKey is a file
	actualPublicKey := defaultPublicKey
	if len(publicKey) == 1 && len(publicKey[0]) != 0 {
		actualPublicKey = publicKey[0]
	}

	block, _ := pem.Decode(actualPublicKey)
	if block == nil {
		return nil, crypto.ErrInvalidPublicKey
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, input)
}

// Decrypt from rsa string
func Decrypt(ciphertext []byte, privateKey ...[]byte) ([]byte, error) {
	// decrypt
	actualPrivateKey := defaultPrivateKey
	if len(privateKey) == 1 && len(privateKey[0]) != 0 {
		actualPrivateKey = privateKey[0]
	}
	block, _ := pem.Decode(actualPrivateKey)
	if block == nil {
		return nil, crypto.ErrInvalidPrivateKey
	}
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, pk, ciphertext)
}

// EncryptBase64 encrypt to base64 string
func EncryptBase64(input []byte, publicKey ...[]byte) ([]byte, error) {
	data, err := Encrypt([]byte(input), publicKey...)
	data = base64.Encode(data)
	return data, err
}

// DecryptBase64 decrypt from base64 string
func DecryptBase64(input []byte, privateKey ...[]byte) ([]byte, error) {
	ciphertext, err := base64.Decode(input)
	if err == nil {
		data, err := Decrypt(ciphertext, privateKey...)
		return data, err
	}
	return nil, crypto.ErrInvalidInput
}
