package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secret string
}

func NewJWTService(secret string) *jwtService {
	return &jwtService{secret: secret}
}

func (s *jwtService) Generate(userID int, name string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"name":    name,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *jwtService) Validate(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		return 0, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrSignatureInvalid
	}

	userID := int(claims["user_id"].(float64))
	return userID, nil
}
