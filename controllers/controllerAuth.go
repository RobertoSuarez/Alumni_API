package controllers

import (
	"net/http"

	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type ControllerAuth struct{}

func NewControllerAuth() *ControllerAuth {
	return &ControllerAuth{}
}

func (c *ControllerAuth) ConfigPath(router fiber.Router) {

	router.Post("/login", c.LoginHandler)

	router.Get("/vista", ValidarJWT, func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*models.Claim)

		return c.SendString("exito perfectamente autenticado " + claims.Name)
	})
}

func (auth *ControllerAuth) LoginHandler(c *fiber.Ctx) error {
	var login models.Login

	err := c.BodyParser(&login)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al convertir los datos"})
	}

	token, err := GenerarJWT(login)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al contruir el token"})
	}

	respuestaLogin := models.RespuestaLogin{
		Token: token,
		Perfil: models.PerfilUser{
			Name:  login.Username,
			Email: login.Username,
		},
	}

	return c.Status(http.StatusOK).JSON(&respuestaLogin)
}
