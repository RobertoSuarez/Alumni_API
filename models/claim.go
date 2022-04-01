package models

import "github.com/golang-jwt/jwt"

type Claim struct {
	Email  string  `json:"email"`
	IdUser uint64  `json:"id"`
	Grupos []Grupo `json:"grupos"`
	jwt.StandardClaims
}
