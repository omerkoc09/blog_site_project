// service/token.go
package service

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService struct {
	secretKey string
}

type TokenParser interface {
	ParseToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID uint `json:"uid"`
	jwt.StandardClaims
}

func NewTokenService(secretKey string) *TokenService {
	return &TokenService{
		secretKey: secretKey,
	}
}

func (s *TokenService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
