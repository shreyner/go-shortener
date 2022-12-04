package sign

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
)

// StringSign base struct for crypto string
//
//	stringSign, err := sign.NewStringSign([]byte("key"))
//	stringSign.Encrypt("123") // 99d3f20d20699649c2f1620c3276cf3585d9608493
//	stringSign.Decrypt("99d3f20d20699649c2f1620c3276cf3585d9608493") // 123
type StringSign struct {
	aesGCM cipher.AEAD
	nonce  []byte
}

// NewStringSign constructor for StringSign
func NewStringSign(key []byte) (*StringSign, error) {
	sh := sha256.New()
	sh.Write(key)

	keyHash := sh.Sum(nil)

	aesBlock, err := aes.NewCipher(keyHash)

	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonce := keyHash[len(keyHash)-aesGCM.NonceSize():]

	stringSign := StringSign{
		aesGCM: aesGCM,
		nonce:  nonce,
	}

	return &stringSign, nil
}

// Encrypt string
func (s *StringSign) Encrypt(str string) string {
	dst := s.aesGCM.Seal(nil, s.nonce, []byte(str), nil)

	return hex.EncodeToString(dst)
}

// Decrypt string
func (s *StringSign) Decrypt(str string) (string, error) {
	v, err := hex.DecodeString(str)

	if err != nil {
		return "", err
	}

	res, err := s.aesGCM.Open(nil, s.nonce, v, nil)

	if err != nil {
		return "", err
	}

	return string(res), nil
}
