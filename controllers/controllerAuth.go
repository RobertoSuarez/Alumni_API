package controllers

import (
	"net/http"

	"github.com/RobertoSuarez/apialumni/database"
	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type ControllerAuth struct{}

func NewControllerAuth() *ControllerAuth {
	return &ControllerAuth{}
}

func (c *ControllerAuth) ConfigPath(router fiber.Router) {

	router.Post("/login", c.LoginHandler)

	router.Get("/users", ValidarJWT, func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*models.Claim)
		_ = claims
		//fmt.Println(claims)
		users := []models.Usuario{}

		database.Database.Find(&users)

		return c.JSON(users)
	})
}

func (auth *ControllerAuth) LoginHandler(c *fiber.Ctx) error {
	var login models.Login
	var usuario models.Usuario

	err := c.BodyParser(&login)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al convertir los datos"})
	}

	result := database.Database.Where("email = ? AND password = ?", login.Username, login.Password).First(&usuario)

	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "El usuario o contraseña esta incorrecto"})
	}

	token, err := GenerarJWT(usuario)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al construir el token"})
	}
	usuario.Password = ""

	respuestaLogin := models.RespuestaLogin{
		Token:   token,
		Usuario: usuario,
	}

	return c.Status(http.StatusOK).JSON(&respuestaLogin)
}
