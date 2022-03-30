package controllers

import (
	"fmt"
	"strconv"

	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type Empresa struct{}

func NewEmpresa() *Empresa {
	return &Empresa{}
}

func (empresa *Empresa) ConfigPath(router *fiber.App) *fiber.App {
	router.Get("/", empresa.ObtenerEmpresas)
	router.Post("/", empresa.CrearEmpresa)
	router.Put("/:id", empresa.Actualizar)

	return router
}

// Envia todas la empresas
func (Empresa) ObtenerEmpresas(c *fiber.Ctx) error {
	empresas, err := models.Empresa{}.ObtenerEmpresas()
	if err != nil {
		return c.Status(400).JSON(err)
	}

	return c.JSON(empresas)
}

// Crear una empresa
func (Empresa) CrearEmpresa(c *fiber.Ctx) error {
	empresa := models.Empresa{}

	err := c.BodyParser(&empresa)
	if err != nil {
		return c.Status(400).JSON(err)
	}

	err = empresa.CrearEmpresa()
	if err != nil {
		fmt.Println("Error en la base de datos: ", err)
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON(empresa)
}

// Actualizar una emprea
func (Empresa) Actualizar(c *fiber.Ctx) error {
	idempresa := c.Params("id")
	ID, err := strconv.ParseInt(idempresa, 10, 64)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	empresa := models.Empresa{}

	err = c.BodyParser(&empresa)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// Establicemos el id que se pasa en la url
	empresa.ID = uint64(ID)

	err = empresa.Actualizar()
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON(empresa)
}
