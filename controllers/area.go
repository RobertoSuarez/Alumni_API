package controllers

import (
	"strconv"

	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type Area struct{}

func NewControllerArea() *Area {
	return &Area{}
}

func (area *Area) ConfigPath(router *fiber.App) *fiber.App {

	router.Get("/", area.ObtenerAreas)
	router.Get("/:id/subareas", area.ObtenerSubareas)
	return router
}

func (Area) ObtenerAreas(c *fiber.Ctx) error {

	areas, err := models.Area{}.ObtenerAreas()
	if err != nil {
		return c.Status(400).SendString("Error en la DB")
	}

	return c.JSON(areas)
}

func (Area) ObtenerSubareas(c *fiber.Ctx) error {
	idarea := c.Params("id")
	ID, err := strconv.ParseInt(idarea, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	subareas, err := models.Subarea{}.ObtenerSubareas(uint64(ID))
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.JSON(subareas)
}
