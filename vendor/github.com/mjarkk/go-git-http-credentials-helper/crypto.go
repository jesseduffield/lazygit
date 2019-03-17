package gitcredentialhelper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
	mathRand "math/rand"
)

// encrypt text using the key
func encrypt(text, key string) ([]byte, error) {
	if key == "" {
		return nil, errors.New("The key can't be empty")
	}

	gcm, err := createGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, []byte(text), nil), nil
}

// decrypt a string output from key
func decrypt(data []byte, key string) (string, error) {
	gcm, err := createGCM(key)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	out, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func createGCM(key string) (cipher.AEAD, error) {
	sha256Key := sha256.Sum256([]byte(key))
	aesCipher, err := aes.NewCipher(sha256Key[:])
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(aesCipher)
}

// encryptWithPublicKey encrypts a string using a public key
func encryptWithPublicKey(publicKey, data string) ([]byte, error) {
	pubBlock, _ := pem.Decode([]byte(publicKey))
	pub, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, []byte(data), []byte{})
}

// decryptWithPrivateKey decrypts data to a string using a private key
func decryptWithPrivateKey(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	output, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, []byte{})
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// generateKeyPair generates a RSA key pair and returns the privatekey struct
// the public key is send back as string becuase it's usualy directly send await to somewhere else
func generateKeyPair() (*rsa.PrivateKey, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, "", err
	}

	publicKey := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)})

	return privateKey, string(publicKey), nil
}

// EncryptedMessage is what a encrypted message looks like
// The key is the encryption key encrypted using the public key of the sender
// The Data is the actual data
type EncryptedMessage struct {
	EncryptedKey []byte
	Data         []byte
}

func decryptMessage(privateKey *rsa.PrivateKey, data EncryptedMessage) (string, error) {
	pass, err := decryptWithPrivateKey(privateKey, data.EncryptedKey)
	if err != nil {
		return "", err
	}

	return decrypt(data.Data, pass)
}

func encryptMessage(publicKey, message string) (EncryptedMessage, error) {
	encryptedMessage := EncryptedMessage{}

	pass, err := randomString(50)
	if err != nil {
		return encryptedMessage, err
	}

	encryptedMessage.Data, err = encrypt(message, pass)
	if err != nil {
		return encryptedMessage, err
	}

	encryptedMessage.EncryptedKey, err = encryptWithPublicKey(publicKey, pass)
	if err != nil {
		return encryptedMessage, err
	}

	return encryptedMessage, nil
}

// randomString generates a purly random string with the lenght of n
func randomString(length int) (string, error) {
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(2147483647))
	if err != nil {
		return "", err
	}
	r := mathRand.New(mathRand.NewSource(randomNumber.Int64()))
	possibleLetters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	l := int64(len(possibleLetters))

	toReturn := ""
	for i := 0; i < length; i++ {
		toReturn = toReturn + string(possibleLetters[r.Int63n(l)])
	}

	return toReturn, nil
}
