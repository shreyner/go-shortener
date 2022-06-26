package middlewares

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/shreyner/go-shortener/internal/pkg/random"
)

var (
	lengthUserID  = 5
	authCookieKey = "auth"
	userCtxKey    = "userID"
)

func GetUserIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(userCtxKey).(string)
	return v
}

func AuthHandler(key []byte) func(next http.Handler) http.Handler {
	sh := sha256.New()
	sh.Write(key)

	keyHash := sh.Sum(nil)

	aesBlock, err := aes.NewCipher(keyHash)
	if err != nil {
		log.Fatalln(err)
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Fatalln(err)
	}

	nonce := keyHash[len(keyHash)-aesGCM.NonceSize():]

	parseCookie := func(r *http.Request) (string, error) {
		authCookie, err := r.Cookie(authCookieKey)

		if err != nil {
			return "", err
		}

		v, err := hex.DecodeString(authCookie.Value)

		if err != nil {
			return "", err
		}

		userID, err := aesGCM.Open(nil, nonce, v, nil)

		if err != nil {
			return "", err
		}

		return string(userID), nil
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			userID, err := parseCookie(r)

			if err != nil {
				newUserID := generateUserID()

				dst := aesGCM.Seal(nil, nonce, []byte(newUserID), nil)

				authCookie := &http.Cookie{Value: hex.EncodeToString(dst), Name: authCookieKey}

				http.SetCookie(rw, authCookie)

				ctx := context.WithValue(r.Context(), userCtxKey, newUserID)

				next.ServeHTTP(rw, r.WithContext(ctx))

				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, userID)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func generateUserID() string {
	return random.RandSeq(lengthUserID)
}
