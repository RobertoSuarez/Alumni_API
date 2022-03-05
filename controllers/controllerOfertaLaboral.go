package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RobertoSuarez/apialumni/database"
	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	router.Get("/:idempleo", cofertas.GetOfertaByID)
}

// Endpoint empleos
// este controlador tendra la tarea de tomar los datos de entrada tales como:
// - Lista de categorias
// - Lista de ciudades
// - Palabras claves u oraciones
// Convertirlos a un struct y buscar los registros que coinsidad en la base de datos
// luego enviar los registros de empleos al cliente en formato json.
func (cofetas *ControllerOfertaLaboral) ObtenerOfetasLaborales(c *fiber.Ctx) error {
	ofertas := []*models.Empleo{}

	query := models.QueryEmpleo{}

	err := c.QueryParser(&query)
	if err != nil {
		log.Println(err)
	}

	condiciones := make(map[string]interface{})

	if len(query.Areas) > 0 && len(query.Areas[0]) > 0 {
		condiciones["area"] = query.Areas
	}

	// Cuando se estable la variable en la url, se estable un valor como cadena vacia
	// con lo cual, se debe revisar la primera posición
	if len(query.Ciudades) > 0 && len(query.Ciudades[0]) > 0 {
		condiciones["ciudad"] = query.Ciudades
	}

	consulta := database.Database.Where(condiciones)

	// si el usuario no ingreso ningun texto, no se considara la busquedad por el titulo
	if len(query.Busquedad) > 0 {
		consulta = consulta.Where("titulo like ?", "%"+query.Busquedad+"%")
	}

	result := consulta.Preload("Usuario", func(tx *gorm.DB) *gorm.DB {
		return tx.Select([]string{
			"ID",
			"IdentificacionTipo",
			"NumeroIdentificacion",
			"Nombres",
			"Apellidos",
			"Email",
			//"Password",
			"Nacimiento",
			"Whatsapp",
			"RoleCuenta",
		})
	}).Find(&ofertas)

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

	oferta := models.Empleo{}

	err := c.BodyParser(&oferta)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al convertir los datos"})
	}

	oferta.UsuarioID = claims.IdUser

	result := database.Database.Create(&oferta)

	if result.Error != nil {
		log.Println(result.Error)
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudo registrar la oferta laboral"})
	}

	//database.Database.Where("id = ?", oferta.ID).First(&oferta)

	return c.Status(http.StatusOK).JSON(oferta)

}

// endpoint getEmpleo debe retornar el registro del empleo en base al id enviado por el cliente
func (ofertas *ControllerOfertaLaboral) GetOfertaByID(c *fiber.Ctx) error {

	empleo := models.Empleo{}

	idusuario := c.Params("idempleo")

	//fmt.Println(idusuario)

	result := database.Database.Where("id = ?", idusuario).Find(&empleo)
	//fmt.Println(result.Error)
	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se encontro el empleo"})
	}

	if result.Statement.RowsAffected < 1 {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No existe ese registro"})
	}

	return c.JSON(empleo)
}
