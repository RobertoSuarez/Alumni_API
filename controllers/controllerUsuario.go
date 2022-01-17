package controllers

import (
	"net/http"

	"github.com/RobertoSuarez/apialumni/database"
	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type ControllerUsuario struct{}

func NewControllerUsuario() *ControllerUsuario {
	return &ControllerUsuario{}
}

func (cuser *ControllerUsuario) ConfigPath(router fiber.Router) {

	router.Get("/", cuser.GetUsuarioHandler)
}

// Retorna el usuario que se autentica con el token
func (cuser *ControllerUsuario) GetUsuarioHandler(c *fiber.Ctx) error {
	usuarios := []models.Usuario{}
	database.Database.Find(&usuarios)

	return c.Status(http.StatusOK).JSON(usuarios)
}
