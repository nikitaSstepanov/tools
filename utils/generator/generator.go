package generator

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	mrand "math/rand"
	"strings"
)

const (
	charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lenOfSymbols = 62
	lenOfNumbs   = 10
)

var (
	symbols = []rune(charset)
	numbs   = []rune("1234567890")
)

func GetRandomNum(length int) string {
	res := make([]rune, length)

	for i := range res {
		res[i] = numbs[mrand.Intn(lenOfNumbs)]
	}

	return string(res)
}

func GetSecret(length int) (string, error) {
	b := make([]byte, length)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}

	key := hex.EncodeToString(b)

	var sb strings.Builder
	for _, r := range key {
		if strings.ContainsRune(charset, r) {
			sb.WriteRune(r)
		}
	}

	for len(sb.String()) < length {
		sb.WriteRune(symbols[mrand.Intn(lenOfSymbols)])
	}

	return sb.String()[:length], nil
}
