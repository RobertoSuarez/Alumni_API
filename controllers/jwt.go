package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var Miclave = []byte("alumniuteq")

func GenerarJWT(u *models.Usuario) (string, error) {

	payload := jwt.MapClaims{
		"email":  u.Email,
		"iduser": u.ID,
		//"tipoUsuario": u.TipoUsuario.Tipo,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(Miclave)
}

func ValidarJWT(c *fiber.Ctx) error {

	var claims models.Claim

	tokenstr := c.Get("Authorization", "")
	if tokenstr == "" {
		return c.Status(http.StatusBadRequest).SendString("Para tener acceso se requiere del token")
	}

	SplitToken := strings.Split(tokenstr, "Bearer")
	if len(SplitToken) != 2 {
		return c.Status(http.StatusBadRequest).SendString("Error en el formato del token")
	}

	tokenstr = strings.TrimSpace(SplitToken[1])

	// Recuperamos los claims del token
	token, err := jwt.ParseWithClaims(tokenstr, &claims, func(t *jwt.Token) (interface{}, error) {
		return Miclave, nil
	})
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Fallo en el token")
	}

	if !token.Valid {
		return c.Status(http.StatusBadRequest).SendString("El token no es validdo")
	}

	c.Locals("claims", &claims)

	fmt.Println(claims)
	return c.Next()
}
