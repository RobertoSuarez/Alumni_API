package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var Miclave = []byte("alumniuteq")

func GenerarJWT(u models.Usuario) (string, error) {

	payload := jwt.MapClaims{
		"email":  u.Email,
		"id":     u.ID,
		"exp":    time.Now().Add(time.Hour * 8).Unix(),
		"grupos": u.Grupos,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(Miclave)
}

//Middleware
func ValidarJWT(c *fiber.Ctx) error {

	claims := models.Claim{}

	tokenstr := c.Get("Authorization", "")
	if tokenstr == "" {
		return c.Status(http.StatusUnauthorized).SendString("Para tener acceso se requiere del token")
	}

	SplitToken := strings.Split(tokenstr, "Bearer")
	if len(SplitToken) != 2 {
		return c.Status(http.StatusUnauthorized).SendString("Error en el formato del token")
	}

	tokenstr = strings.TrimSpace(SplitToken[1])

	// Recuperamos los claims del token
	token, err := jwt.ParseWithClaims(tokenstr, &claims, func(t *jwt.Token) (interface{}, error) {
		return Miclave, nil
	})
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("Fallo en el token")
	}

	if !token.Valid {
		return c.Status(http.StatusUnauthorized).SendString("El token no es validdo")
	}

	c.Locals("claims", &claims)

	//fmt.Println(claims)
	return c.Next()
}

//Esta funcion retorna una función con los grupos permitidos
func gruposPermitios(grupos []string) func(*fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*models.Claim)

		for _, permitido := range grupos {
			for _, pertenece := range claims.Grupos {
				if permitido == pertenece.Name {
					//fmt.Println("si tienes acceso")
					return c.Next()
				}
			}
		}
		//fmt.Println("No tienes acceso")
		// el usuario no tiene siertos privilegios
		return c.Status(http.StatusForbidden).JSON(models.NewError("no tienes los privilegios requeridos"))
	}
}
