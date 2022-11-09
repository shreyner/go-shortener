// Package random include helpers for generate random strings by length string
package random

import (
	"crypto/rand"
	"math/big"
)

func cryptoRandomInt(min, max int) (int, error) {
	bg := big.NewInt(int64(max - min + 1))

	n, err := rand.Int(rand.Reader, bg)

	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + min, nil
}

// letters include all characters by generate
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// RandSeq random string by based letters
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		rn, err := cryptoRandomInt(0, len(letters)-1)

		if err != nil {
			panic(err)
		}

		b[i] = letters[rn]
	}
	return string(b)
}
