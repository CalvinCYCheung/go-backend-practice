package tokenservice

import (
	"time"

	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(userId string) (*string, error)
	ValidateToken(token string) (*string, error)
}

type TokenServiceImpl struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewTokenService(secretKey string, tokenDuration time.Duration) TokenService {
	// TODO: Get latest private key from S3 bucket
	return &TokenServiceImpl{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (t *TokenServiceImpl) GenerateToken(userId string) (*string, error) {
	// TODO: Generate token with private key
	var privateKey = &rsa.PrivateKey{}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(t.tokenDuration).Unix(),
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (t *TokenServiceImpl) ValidateToken(token string) (*string, error) {
	// TODO: Get public key from S3 bucket , return error if the key is not found
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return claims["sub"].(*string), nil
}
