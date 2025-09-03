package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AdminID string `json:"AdminID"`
	jwt.RegisteredClaims
}
