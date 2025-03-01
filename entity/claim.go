package entity

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email string
	jwt.RegisteredClaims
}
