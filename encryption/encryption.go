package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// RSADecrypt implement decryption with rsa.
func RSADecrypt(cipherText string, privKey rsa.PrivateKey) (string, error) {
	ct, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	plaintext, err := rsa.DecryptOAEP(sha256.New(),
		rand.Reader, &privKey, ct, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// GetPasswordHASH return hashed password.
func GetPasswordHASH(passwd []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(passwd,
		bcrypt.DefaultCost)
}

// CheckPasswordHASH check password with hash.
func CheckPasswordHASH(passwd, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash,
		passwd)
}
