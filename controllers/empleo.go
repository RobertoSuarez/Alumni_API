package controllers

import (
	"strconv"

	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type Empleo struct {
}

func NewEmpleo() *Empleo {
	return &Empleo{}
}

func (e *Empleo) ConfigPath(app *fiber.App) *fiber.App {

	app.Get("/", e.ListarEmpleos)
	app.Post("/", e.Crear)
	app.Put("/:id", e.Actualizar)
	return app
}

// obtener todos los empleos registrador
func (e *Empleo) ListarEmpleos(c *fiber.Ctx) error {

	empleos, err := models.Empleo{}.ObtenerTodos()
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.JSON(empleos)
}

// Publicar Empleo
func (e *Empleo) Crear(c *fiber.Ctx) error {
	empleo := models.Empleo{}

	err := c.BodyParser(&empleo)
	if err != nil {
		return c.Status(400).SendString("Error al convertir los datos")
	}

	err = empleo.Crear()
	if err != nil {
		return c.Status(400).SendString("Error: " + err.Error())
	}

	return c.JSON(empleo)
}

// Actualizar empleo
func (e *Empleo) Actualizar(c *fiber.Ctx) error {
	empleo := models.Empleo{}
	idempleo := c.Params("id")
	ID, err := strconv.ParseInt(idempleo, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	if err = c.BodyParser(&empleo); err != nil {
		return c.Status(400).SendString("Error al convertir los datos")
	}

	empleo.ID = uint64(ID)

	err = empleo.Actualizar()
	if err != nil {
		return c.Status(400).SendString("Error: " + err.Error())
	}

	return c.JSON(empleo)
}
