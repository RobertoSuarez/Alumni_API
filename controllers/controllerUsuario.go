package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RobertoSuarez/apialumni/database"
	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type ControllerUsuario struct{}

func NewControllerUsuario() *ControllerUsuario {
	return &ControllerUsuario{}
}

func (cuser *ControllerUsuario) ConfigPath(router fiber.Router) {

	router.Get("/tipos", cuser.GetTiposUsuario)
	router.Get("/", cuser.GetUsuarioHandler)
	router.Post("/", cuser.CrearUsuarioHandler)

}

// Retorna el usuario que se autentica con el token
func (cuser *ControllerUsuario) GetUsuarioHandler(c *fiber.Ctx) error {
	usuarios := []*models.Usuario{}
	database.Database.Select(
		[]string{
			"ID",
			"IdentificacionTipo",
			"NumeroIdentificacion",
			"Nombres",
			"Apellidos",
			"Email",
			//"Password",
			"Nacimiento",
			"Whatsapp",
			"Administrador",
		}).Find(&usuarios)
	//database.Database.Preload("Admin").Preload("Alumni").Preload("TipoUsuario").Find(&usuarios)

	// for _, usuario := range usuarios {
	// 	// if usuario.TipoUsuario.Tipo == "admin" {
	// 	// 	usuario.Alumni = nil
	// 	// }

	// 	// if usuario.TipoUsuario.Tipo == "alumni" {
	// 	// 	usuario.Admin = nil
	// 	// }
	// 	ClearTiposUsuarios(usuario)
	// }

	return c.Status(http.StatusOK).JSON(usuarios)
}

// ClearTiposUsuarios establece el valor nil a los tipos de usuarios que no se define el usuario
// func ClearTiposUsuarios(usuario *models.Usuario) {
// 	if usuario.TipoUsuario.Tipo == "admin" {
// 		usuario.Alumni = nil
// 	}

// 	if usuario.TipoUsuario.Tipo == "alumni" {
// 		usuario.Admin = nil
// 	}
// }

func (cuser *ControllerUsuario) CrearUsuarioHandler(c *fiber.Ctx) error {
	tipos := models.ListTipoUsuarios{}
	usuario := models.Usuario{}

	err := c.BodyParser(&usuario)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudieron convertir los datos"})
	}

	// Verificamos si ya existe un usuario registrado con el email.
	user := models.Usuario{}
	result := database.Database.Where("email = ?", usuario.Email).First(&user)
	if result.RowsAffected > 0 {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "El email ya existe"})
	}

	// verificar las estructura de alumni admin
	// fmt.Println("Objeto usuario completo: ", usuario)
	// if usuario.Admin == nil {
	// 	usuario.Admin = new(models.Admin).SetNil()
	// }
	// if usuario.Alumni == nil {
	// 	usuario.Alumni = new(models.Alumni).SetNil()
	// }

	fmt.Println("Objeto usuario completo con objetos en nil: ", usuario)

	database.Database.Find(&tipos)
	fmt.Println(tipos)
	//fmt.Println(tipos.GetID(usuario.TipoUsuario.Tipo))
	// //fmt.Println(tipos)
	// tipo, err := convertirID_TipoUsuario(tipos, usuario.TipoUsuarioID)
	// if err != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "El tipo de usuario que enviastes no existe"})
	// }

	// revisamos que tipos de usuario se va a registrar
	// fmt.Println(usuario.TipoUsuario)
	// switch usuario.TipoUsuario.Tipo {
	// case "admin":

	// 	usuario.Admin.Estado.Usando = true
	// 	usuario.Alumni.SetNil()
	// case "alumni":
	// 	usuario.Alumni.Estado.Usando = true
	// 	usuario.Admin.SetNil()
	// default:
	// 	fmt.Println("No esta tomando ningun caso")
	// 	return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "El tipo de usuario no existe"})
	// }

	// usuario.TipoUsuario.ID, _ = tipos.GetID(usuario.TipoUsuario.Tipo) // Establecemos el id del tipo de usuario.

	// creamos el usuario
	result = database.Database.Create(&usuario)
	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al registrar el usuario"})
	}

	return c.SendStatus(http.StatusOK)
}

// Convierte el id del tipos de usuario en el nombre del tipo reguistrado en la db
func convertirID_TipoUsuario(tipos []models.TipoUsuario, id uint) (string, error) {
	for _, v := range tipos {
		if v.ID == id {
			return v.Tipo, nil
		}
	}
	return "", errors.New("no existe el tipo")
}

// endpoint para los tipos de usuarios
func (cuser *ControllerUsuario) GetTiposUsuario(c *fiber.Ctx) error {
	tipos := []models.TipoUsuario{}

	database.Database.Find(&tipos)

	return c.JSON(tipos)
}
