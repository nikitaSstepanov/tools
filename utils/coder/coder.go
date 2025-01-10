package coder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Secret   string `env:"ENCRYPT_SECRET" env-default:""`
	HashCost int    `env:"HASH_COST" env-default:"10"`
}

type Coder struct {
	secret string
	cost   int
}

func New(cfg *Config) *Coder {
	return &Coder{
		secret: cfg.Secret,
		cost:   cfg.HashCost,
	}
}

func (c *Coder) Encrypt(text string) string {
	aes, err := aes.NewCipher([]byte(c.secret))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)

	return hex.EncodeToString(ciphertext)
}

func (c *Coder) Decrypt(text string) (string, error) {
	b, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}

	text = string(b)

	aes, err := aes.NewCipher([]byte(c.secret))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := text[:nonceSize], text[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (c *Coder) Hash(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), c.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (c *Coder) CompareHash(hash string, text string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
}
