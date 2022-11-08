package middlewares

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/pkg/random"
)

const userCtxKey = iota

var (
	lengthUserID  = 5
	authCookieKey = "auth"
)

// GetUserIDFromCtx return user ID from contect
func GetUserIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(userCtxKey).(string)
	return v
}

// AuthHandler for auth users and create if not found auth cookies
func AuthHandler(log *zap.Logger, key []byte) func(next http.Handler) http.Handler {
	sh := sha256.New()
	sh.Write(key)

	keyHash := sh.Sum(nil)

	aesBlock, err := aes.NewCipher(keyHash)
	if err != nil {
		// TODO: Убрать Fatalln. Обычный error. Сделать выбрасывание http ошибки
		log.Fatal("error", zap.Error(err))
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		// TODO: Убрать Fatalln. Обычный error. Сделать выбрасывание http ошибки
		// TODO: Пробежаться и посомтреть по коду
		log.Fatal("error", zap.Error(err))
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
