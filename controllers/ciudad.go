package controllers

import (
	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type Ciudad struct{}

func NewControllerCiudad() *Ciudad {
	return &Ciudad{}
}

func (ciudad *Ciudad) ConfigPath(router *fiber.App) *fiber.App {

	router.Get("/", ciudad.ObtenerTodas)
	return router
}

func (Ciudad) ObtenerTodas(c *fiber.Ctx) error {

	ciudades, err := models.Ciudad{}.ObtenerTodas()
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	return c.JSON(ciudades)
}
