package random

import (
	"crypto/rand"
	"math/big"
)

func CryptoGenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)

	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func CryptoRandomInt(min, max int) (int, error) {
	bg := big.NewInt(int64(max - min + 1))

	n, err := rand.Int(rand.Reader, bg)

	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + min, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		rn, err := CryptoRandomInt(0, len(letters)-1)

		if err != nil {
			panic(err)
		}

		b[i] = letters[rn]
	}
	return string(b)
}
