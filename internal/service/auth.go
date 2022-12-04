package service

import (
	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/pkg/random"
	"github.com/shreyner/go-shortener/internal/pkg/sign"
)

var lengthUserID = 5

// AuthService include base method for authantification in service
type AuthService struct {
	log        *zap.Logger
	stringSign *sign.StringSign
}

// NewAuthService constructor
func NewAuthService(log *zap.Logger, signKey []byte) (*AuthService, error) {
	stringSign, err := sign.NewStringSign(signKey)

	if err != nil {
		return nil, err
	}

	authService := AuthService{
		log:        log,
		stringSign: stringSign,
	}

	return &authService, nil
}

// GenerateUserID Generate userID
//
// TODO: Move to UserService
func (s *AuthService) GenerateUserID() string {
	return generateUserID()
}

// CreateToken by userID
func (s *AuthService) CreateToken(userID string) string {
	return s.stringSign.Encrypt(userID)
}

// GetUserIDFromToken by token
func (s *AuthService) GetUserIDFromToken(token string) (string, error) {
	return s.stringSign.Decrypt(token)
}

func generateUserID() string {
	return random.RandSeq(lengthUserID)
}
