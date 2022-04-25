package controllers

import (
	"strconv"

	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type Provincia struct{}

func NewControllerProvincia() *Provincia {
	return &Provincia{}
}

func (provincia *Provincia) ConfigPath(router *fiber.App) *fiber.App {
	router.Get("/", provincia.ObtenerTodas)
	router.Get("/:id/ciudades", provincia.ObtenerCiudades)
	return router
}

// Lista a las provincias
func (Provincia) ObtenerTodas(c *fiber.Ctx) error {

	provincias, err := models.Provincia{}.ObtenerTodas()
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.JSON(provincias)
}

// Obtener las ciudades por el id de la provincia
func (Provincia) ObtenerCiudades(c *fiber.Ctx) error {

	ID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	ciudades, err := models.Provincia{ID: uint64(ID)}.ObtenerCiudades()
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.JSON(ciudades)
}
