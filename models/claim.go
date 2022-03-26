package models

import "github.com/golang-jwt/jwt"

type Claim struct {
	Email       string `json:"email"`
	TipoUsuario string `json:"tipoUsuario"`
	IdUser      uint64 `json:"iduser"`
	jwt.StandardClaims
}
