package rsa

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/crypto"
	"testing"
)

var invalidPrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MCsCAQACBQDJaQRdAgMBAAECBEPxlU0CAwDt4wIDANi/AgJC0QICBc0CAkg4
-----END RSA PRIVATE KEY-----
`)

var invalidPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MCAwDQYJKoZIhvcNAQEBBQADDwAwDAIFAMlpBF0CAwEAAQ==
-----END PUBLIC KEY-----
`)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestRsa(t *testing.T) {
	src := []byte("hello")
	data, _ := Encrypt([]byte(src))
	decrypted, err := Decrypt(data)
	assert.Equal(t, nil, err)
	log.Debugf("encrypted: %v, decrypted: %v, org: %v", string(data), string(decrypted), string(src))
	assert.Equal(t, src, decrypted)
}

func TestRsaBase64(t *testing.T) {
	src := []byte("hello")
	data, _ := EncryptBase64([]byte(src))
	decrypted, err := DecryptBase64(data)
	assert.Equal(t, nil, err)
	log.Debugf("encrypted: %v, decrypted: %v, org: %v", string(data), string(decrypted), string(src))
	assert.Equal(t, src, decrypted)
}

func TestExeptions(t *testing.T) {
	t.Run("should report error with invalid public key", func(t *testing.T) {
		src := []byte("hello")
		_, err := Encrypt([]byte(src), []byte("invalid-key"))
		assert.Equal(t, crypto.ErrInvalidPublicKey, err)
	})

	t.Run("should report error with invalid public key", func(t *testing.T) {
		src := []byte("hello")
		_, err := Encrypt([]byte(src), invalidPrivateKey)
		assert.Contains(t, err.Error(), "failed")
	})

	t.Run("should report error with invalid public key", func(t *testing.T) {
		src := []byte("hello")
		data, _ := Encrypt([]byte(src))
		_, err := Decrypt(data, invalidPublicKey)
		assert.Contains(t, err.Error(), "tags don't match")
	})

	t.Run("should report error with invalid private key", func(t *testing.T) {
		src := []byte("hello")
		data, _ := Encrypt([]byte(src))
		_, err := Decrypt(data, []byte("invalid-key"))
		assert.Equal(t, crypto.ErrInvalidPrivateKey, err)
	})

	t.Run("should report error with invalid base64 string", func(t *testing.T) {
		src := []byte("invalid-base64")
		_, err := DecryptBase64([]byte(src))
		assert.Equal(t, crypto.ErrInvalidInput, err)
	})
}
