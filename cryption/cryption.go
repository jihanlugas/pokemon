package cryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"gopokemon/config"
	"io"

	"golang.org/x/crypto/bcrypt"
)

var aesCipher cipher.Block

// gcm or Galois/Counter Mode, is a mode of operation
// for symmetric key cryptographic block ciphers
// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
var gcm cipher.AEAD
var nonceSize int
var decrError error

func init() {
	var err error

	// generate a new aes cipher using our 32 byte long key
	aesCipher, err = aes.NewCipher(config.CryptoKey)
	if err != nil {
		panic(err)
	}

	gcm, err = cipher.NewGCM(aesCipher)
	if err != nil {
		panic(err)
	}

	nonceSize = gcm.NonceSize()

	decrError = errors.New("failed authentication")
}

func EncryptAES64(payload []byte) ([]byte, error) {
	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, nonceSize)
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	encBytes := gcm.Seal(nonce, nonce, payload, nil)
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(encBytes)))
	base64.StdEncoding.Encode(b64, encBytes)

	return b64, nil
}

func DecryptAES64(b64 []byte) ([]byte, error) {
	decodedBytes := make([]byte, base64.StdEncoding.DecodedLen(len(b64)))
	n, err := base64.StdEncoding.Decode(decodedBytes, b64)
	if err != nil {
		return nil, err
	}
	ciphertext := decodedBytes[:n]

	if len(ciphertext) < nonceSize {
		return nil, decrError
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
