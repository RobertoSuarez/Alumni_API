package controllers

import (
	"net/http"

	"github.com/RobertoSuarez/apialumni/database"
	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type ControllerOfertaLaboral struct {
}

func NewControllerOfertaLaboral() *ControllerOfertaLaboral {
	return &ControllerOfertaLaboral{}
}

func (cofertas *ControllerOfertaLaboral) ConfigPath(router fiber.Router) {

	// Definimos las rutas.
	router.Get("/", cofertas.ObtenerOfetasLaborales)
	router.Post("/", cofertas.CrearOfertaLaboral)
}

func (cofetas *ControllerOfertaLaboral) ObtenerOfetasLaborales(c *fiber.Ctx) error {
	ofertas := []*models.OfertaLaboral{}

	result := database.Database.Preload("Usuario").Find(&ofertas)
	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error en la db"})
	}

	return c.JSON(ofertas)
}

func (ofertas *ControllerOfertaLaboral) CrearOfertaLaboral(c *fiber.Ctx) error {
	oferta := models.OfertaLaboral{}

	err := c.BodyParser(&oferta)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al convertir los datos"})
	}

	result := database.Database.Create(&oferta)

	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudo registrar la oferta laboral"})
	}

	//database.Database.Where("id = ?", oferta.ID).First(&oferta)

	return c.Status(http.StatusOK).JSON(oferta)

}
