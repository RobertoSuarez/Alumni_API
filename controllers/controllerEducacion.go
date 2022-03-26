package controllers

import (
	"log"
	"net/http"

	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type ControllerEducacion struct{}

func NewControllerEducacion() *ControllerEducacion {
	return &ControllerEducacion{}
}

func (educacion *ControllerEducacion) ConfigPath(app *fiber.App) *fiber.App {

	app.Get("/:iduser", educacion.GetEducacionHandler)
	app.Post("/", ValidarJWT, educacion.CreateEducacionHandler)

	return app
}

// Retorna el usuario que se autentica con el token
func (cuser *ControllerEducacion) GetEducacionHandler(c *fiber.Ctx) error {
	idusuario := c.Params("iduser")
	log.Println(idusuario)
	educacion := []*models.Educacion{}

	// database.Database.Where("usuario_id = ?", idusuario).Find(&educacion)

	return c.Status(http.StatusOK).JSON(educacion)
}

// CreateEducacionHandler crea un registro de donde ha estudiado
func (educacion *ControllerEducacion) CreateEducacionHandler(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)

	edu := models.Educacion{}

	err := c.BodyParser(&edu)
	if err != nil {
		log.Println(err)
	}

	edu.UsuarioID = claims.IdUser

	//result := database.Database.Create(&edu)

	// if result.Error != nil {
	// 	log.Println(result.Error)
	// 	return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudo registrar"})
	// }

	return c.Status(http.StatusOK).JSON(edu)
}
