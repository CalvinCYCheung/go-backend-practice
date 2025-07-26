package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Kid string `json:"kid"`
	jwt.MapClaims
}
