package controllers

import (
	"fmt"
	"log"
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
	router.Post("/", ValidarJWT, cofertas.CrearOfertaLaboral)
}

// Endpoint empleos
// este controlador tendra la tarea de tomar los datos de entrada tales como:
// - Lista de categorias
// - Lista de ciudades
// - Palabras claves u oraciones
// Convertirlos a un struct y buscar los registros que coinsidad en la base de datos
// luego enviar los registros de empleos al cliente en formato json.
func (cofetas *ControllerOfertaLaboral) ObtenerOfetasLaborales(c *fiber.Ctx) error {
	ofertas := []*models.OfertaLaboral{}

	query := models.QueryEmpleo{}

	err := c.BodyParser(&query)
	if err != nil {
		log.Println(err)
	}
	condiciones := make(map[string]interface{})

	if len(query.Area) > 0 {
		condiciones["area"] = query.Area
	}

	if len(query.Ciudades) > 0 {
		condiciones["ciudad"] = query.Ciudades
	}

	consulta := database.Database.Where(condiciones)

	// si el usuario no ingreso ningun texto, no se considara la busquedad por el titulo
	if len(query.Busquedad) > 0 {
		consulta = consulta.Where("titulo like ?", "%"+query.Busquedad+"%")
	}

	result := consulta.Find(&ofertas)

	//result := database.Database.Preload("Usuario").Find(&ofertas)
	// result := database.Database.Preload("Usuario", func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Select([]string{
	// 		"ID",
	// 		"IdentificacionTipo",
	// 		"NumeroIdentificacion",
	// 		"Nombres",
	// 		"Apellidos",
	// 		"Email",
	// 		//"Password",
	// 		"Nacimiento",
	// 		"Whatsapp",
	// 		"RoleCuenta",
	// 	})
	// }).Find(&ofertas)

	if result.Error != nil {
		log.Println(result.Error)
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error en la db"})
	}

	return c.JSON(ofertas)
}

// La oferta laboral se crea con el id usuario que ha iniciado sesión
func (ofertas *ControllerOfertaLaboral) CrearOfertaLaboral(c *fiber.Ctx) error {

	claims := c.Locals("claims").(*models.Claim)

	fmt.Println("Creando oferta laboral, ", claims.IdUser)

	oferta := models.OfertaLaboral{}

	err := c.BodyParser(&oferta)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al convertir los datos"})
	}

	oferta.UsuarioID = claims.IdUser

	result := database.Database.Create(&oferta)

	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudo registrar la oferta laboral"})
	}

	//database.Database.Where("id = ?", oferta.ID).First(&oferta)

	return c.Status(http.StatusOK).JSON(oferta)

}
