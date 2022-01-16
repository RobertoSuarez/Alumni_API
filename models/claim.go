package models

import "github.com/golang-jwt/jwt"

type Claim struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.StandardClaims
}
