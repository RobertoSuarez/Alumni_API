package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/RobertoSuarez/apialumni/awss3"
	"github.com/RobertoSuarez/apialumni/database"
	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ControllerUsuario struct{}

func NewControllerUsuario() *ControllerUsuario {
	return &ControllerUsuario{}
}

func (cuser *ControllerUsuario) ConfigPath(router fiber.Router) {

	router.Post("/", cuser.CrearUsuarioHandler)

	// end point para subir imagen del usuario
	router.Post("/avatar", ValidarJWT, cuser.subirAvatar)
	router.Post("/avataraws", ValidarJWT, cuser.subirAvatarAWS)

	//router.Static("/avatar", "./imgs")
	router.Get("/avatar/:filename", cuser.GetAvatarUsuario)
	router.Get("/avataraws/:filename", cuser.GetAvatarUsuarioAWS)

	router.Get("/", cuser.GetUsuariosHandler)
	router.Get("/tipos", cuser.GetTiposUsuario)
	router.Get("/:iduser", cuser.GetUsuarioByID)

}

// Envia el avatar al cliente
func (cuser *ControllerUsuario) GetAvatarUsuario(c *fiber.Ctx) error {
	filename := c.Params("filename")
	err := c.SendFile("./imgs/" + filename)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No existe ese archivo"})
	}
	return c.SendStatus(http.StatusOK)
}

// GetAvatarUsuario with aws bucket
func (cuser *ControllerUsuario) GetAvatarUsuarioAWS(c *fiber.Ctx) error {
	filename := c.Params("filename")
	fmt.Println("send file image: ", filename)

	resp, err := awss3.GetImage("/fullimages/", filename)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "La imagen no se encontro"})
	}
	defer resp.Body.Close()

	c.Set("Content-Type", *resp.ContentType)
	c.SendStream(resp.Body)
	return c.SendStatus(http.StatusOK)
}

// Retorna el usuario que se autentica con el token
func (cuser *ControllerUsuario) GetUsuariosHandler(c *fiber.Ctx) error {
	usuarios := []*models.Usuario{}
	database.Database.Select(models.UsuarioCamposDB).Find(&usuarios)
	//database.Database.Preload("Admin").Preload("Alumni").Preload("TipoUsuario").Find(&usuarios)

	return c.Status(http.StatusOK).JSON(usuarios)
}

// Recuper de la base de datos, el usuario que se le pasa por id.
func (cuser *ControllerUsuario) GetUsuarioByID(c *fiber.Ctx) error {
	idusuario := c.Params("iduser")

	usuario := models.Usuario{}
	result := database.Database.Where("ID = ?", idusuario).Select(models.UsuarioCamposDB).Find(&usuario)
	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al consultar el usuario"})
	}

	fmt.Println("Usuario recuperado de la base de datos: ", usuario)

	return c.Status(http.StatusOK).JSON(usuario)

}

func (cuser *ControllerUsuario) CrearUsuarioHandler(c *fiber.Ctx) error {

	usuario := models.Usuario{}

	err := c.BodyParser(&usuario)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudieron convertir los datos"})
	}

	if len(usuario.RoleCuenta) < 1 {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "La cuenta debe tener un RoleCuenta"})
	}

	// Verificamos si ya existe un usuario registrado con el email.
	user := models.Usuario{}
	result := database.Database.Where("email = ?", usuario.Email).First(&user)
	if result.RowsAffected > 0 {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "El email ya existe"})
	}

	// creamos el usuario
	result = database.Database.Create(&usuario)
	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al registrar el usuario"})
	}

	return c.SendStatus(http.StatusOK)
}

// endpoint para subir avatar
func (cuser *ControllerUsuario) subirAvatar(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudo procesar esta imagen"})
	}

	// Construimos un nuevo nombre para el archivo que sea unico
	uuid := strings.Replace(uuid.NewString(), "-", "", -1)
	ext := filepath.Ext(file.Filename)
	fileAvatarName := uuid + ext
	fmt.Println(fileAvatarName)

	// Guardamos el archivo y lo registramos en la base de datos.
	err = c.SaveFile(file, fmt.Sprintf("./imgs/%s", fileAvatarName))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudo guardar la imagen"})
	}

	result := database.Database.Model(&models.Usuario{ID: claims.IdUser}).Update("URLAvatar", fileAvatarName)
	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al actualizar el nombre de la imagen"})
	}

	return c.SendStatus(http.StatusOK)
}

// endpoint para subir avatar en el s3 de amazon
func (cuser *ControllerUsuario) subirAvatarAWS(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudo procesar esta imagen"})
	}

	// Construimos un nuevo nombre para el archivo que sea unico
	uuid := strings.Replace(uuid.NewString(), "-", "", -1)
	ext := filepath.Ext(file.Filename)
	fileAvatarName := uuid + ext
	fmt.Println(fileAvatarName)

	// Guardamos el archivo y lo registramos en la base de datos.
	// err = c.SaveFile(file, fmt.Sprintf("./imgs/%s", fileAvatarName))
	// if err != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudo guardar la imagen"})
	// }

	// guardar en aws
	awss3.GuardarImagen("/fullimages/", fileAvatarName, file)

	result := database.Database.Model(&models.Usuario{ID: claims.IdUser}).Update("URLAvatar", fileAvatarName)
	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al actualizar el nombre de la imagen"})
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
