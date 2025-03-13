package models

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Name  string
	Email string
	Admin bool
	jwt.RegisteredClaims
}

func NewClaims(user *User, isAdmin bool, exp time.Time) *Claims {
	return &Claims{
		Name:  user.ID,
		Email: user.Email,
		Admin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
}

func (p *Claims) ToJwtToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, p)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
