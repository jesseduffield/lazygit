package gitcredentialhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	encryptedData, err := encrypt("testString", "test")
	assert.NoError(t, err)
	assert.Greater(t, len(encryptedData), 20)

	output, err := decrypt(encryptedData, "test")
	assert.NoError(t, err)
	assert.Equal(t, output, "testString")

	_, err = decrypt(make([]byte, 32), "test")
	assert.Error(t, err)

	_, err = encrypt("", "")
	assert.Error(t, err)
}

func TestEncryptKeyPair(t *testing.T) {
	longString := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	cases := []string{
		"testInput",
		longString + longString + longString + longString + longString + longString,
	}

	for _, testString := range cases {
		priv, pub, err := generateKeyPair()
		assert.NoError(t, err)

		encryptedData, err := encryptMessage(pub, testString)
		assert.NoError(t, err)

		out, err := decryptMessage(priv, encryptedData)
		assert.NoError(t, err)
		assert.Equal(t, out, testString)
	}
}
