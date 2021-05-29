package aes

import (
	"crypto/aes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/utils/crypto"
	"testing"
)

func TestEncrypt(t *testing.T) {
	originalText := "encrypt this golang text"
	fmt.Println(originalText)
	key := []byte("example key 1234")

	// encrypt value to base64
	t.Run("should encrypt with aes protocol", func(t *testing.T) {
		cryptoText, err := Encrypt(key, originalText)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, "", cryptoText)
	})

	// encrypt base64 crypto to original value
	t.Run("should report to aes.KeySizeError(3) error", func(t *testing.T) {
		_, err := Encrypt([]byte("abc"), originalText)
		assert.Equal(t, aes.KeySizeError(3), err)
	})

	// encrypt base64 crypto to original value
	t.Run("should decrypt with aes protocol", func(t *testing.T) {
		cryptoText, err := Encrypt(key, originalText)
		assert.Equal(t, nil, err)
		text, err := Decrypt(key, cryptoText)
		assert.Equal(t, nil, err)
		assert.Equal(t, originalText, text)
	})

	t.Run("should decrypt with aes protocol", func(t *testing.T) {
		cryptoText, err := Encrypt(key, originalText)
		assert.Equal(t, nil, err)
		_, err = Decrypt(key[4:], cryptoText[40:])
		assert.Equal(t, crypto.ErrCipherTooShort, err)
	})

	// encrypt base64 crypto to original value
	t.Run("should decrypt with aes protocol", func(t *testing.T) {
		cryptoText, err := Encrypt(key, originalText)
		assert.Equal(t, nil, err)
		_, err = Decrypt([]byte("abc"), cryptoText)
		assert.Equal(t, aes.KeySizeError(3), err)
	})
}
