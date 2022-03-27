package controllers

import (
	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type Grupo struct{}

func NewGrupo() *Grupo {
	return &Grupo{}
}

func (g *Grupo) ConfigPath(app *fiber.App) *fiber.App {
	app.Get("/", g.ObtenerGrupos)
	app.Post("/", g.CrearGrupo)
	return app
}

func (g *Grupo) ObtenerGrupos(c *fiber.Ctx) error {
	grupos, err := models.Grupo{}.ObtenerGrupos()
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.JSON(grupos)
}

func (g *Grupo) CrearGrupo(c *fiber.Ctx) error {
	var grupo models.Grupo

	err := c.BodyParser(&grupo)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	err = grupo.Crear()
	if err != nil {
		return c.Status(400).SendString("El nombre del grupo no se puede repetir")
	}

	return c.JSON(grupo)
}
