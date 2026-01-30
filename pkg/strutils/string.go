package strutils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func ExtractNameFromEmail(email string) string {
	name, _, _ := strings.Cut(email, "@")
	return name
}

func RandStr(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[num.Int64()]
	}
	return string(b)
}
