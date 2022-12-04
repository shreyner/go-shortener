package sign

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
)

type StringSign struct {
	aesGCM cipher.AEAD
	nonce  []byte
}

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

func (s *StringSign) Encrypt(str string) string {
	dst := s.aesGCM.Seal(nil, s.nonce, []byte(str), nil)

	return hex.EncodeToString(dst)
}

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
