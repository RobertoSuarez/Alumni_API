package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/RobertoSuarez/apialumni/awss3"
	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Usuario struct{}

func NewControllerUsuario() *Usuario {
	return &Usuario{}
}

func (user *Usuario) ConfigPath(router *fiber.App) *fiber.App {

	router.Get("/", user.ObtenerUsuarios)
	router.Post("/", user.CrearUsuario)
	router.Delete("/", user.EliminarUsuario)
	router.Put("/:id", user.Actualizar)

	router.Post("/confirmar-correo/:id", user.ConfirmarCorreo)

	router.Get("/avatar/:filename", user.GetAvatarUsuario)
	router.Post("/avatar", ValidarJWT, user.subirAvatar)

	router.Get("/avataraws/:filename", user.GetAvatarUsuarioAWS)
	router.Post("/avataraws", ValidarJWT, user.subirAvatarAWS)

	router.Post("/agregar-grupo/:idusuario", user.AgergarGrupo)
	router.Post("/agregar-trabajo/:idusuario", user.AgregarTrabajo)
	router.Get("/:iduser", user.GetUsuarioByID)

	return router
}

func (u *Usuario) EliminarUsuario(c *fiber.Ctx) error {
	type usuariosID struct {
		IDS []uint64 `json:"ids"`
	}
	usuarios := usuariosID{}
	err := c.BodyParser(&usuarios)
	if err != nil {
		return c.Status(400).SendString("error en los datos")
	}
	fmt.Println("usuarios ids: ", usuarios)

	err = models.Usuario{}.Eliminar(usuarios.IDS)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).SendString("no se pudieron eliminar los usuarios")
	}

	return c.SendStatus(http.StatusOK)

}

func (u *Usuario) AgergarGrupo(c *fiber.Ctx) error {
	var grupo models.Grupo
	idusuario := c.Params("idusuario")
	ID, err := strconv.ParseInt(idusuario, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	err = c.BodyParser(&grupo)
	if err != nil {
		return c.Status(400).SendString("Error en el grupo")
	}

	err = models.Usuario{ID: uint64(ID)}.AgregarGrupo(grupo)
	if err != nil {
		return c.Status(400).SendString("No se pudo agregar el grupo")
	}

	return c.SendStatus(http.StatusOK)
}

// Envia el avatar al cliente
func (user *Usuario) GetAvatarUsuario(c *fiber.Ctx) error {
	filename := c.Params("filename")
	err := c.SendFile("./imgs/" + filename)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No existe ese archivo"})
	}
	return c.SendStatus(http.StatusOK)
}

// GetAvatarUsuario with aws bucket
func (user *Usuario) GetAvatarUsuarioAWS(c *fiber.Ctx) error {
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

// Retorna todos los usuarios
func (user *Usuario) ObtenerUsuarios(c *fiber.Ctx) error {
	usuarios, err := models.Usuario{}.GetAll()
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	return c.Status(http.StatusOK).JSON(usuarios)
}

// Recuper de la base de datos, el usuario que se le pasa por id.
func (user *Usuario) GetUsuarioByID(c *fiber.Ctx) error {
	idusuario := c.Params("iduser")
	ID, err := strconv.ParseInt(idusuario, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "El ID no es un entero"})
	}

	usuario, err := models.Usuario{ID: uint64(ID)}.GetUsuarioByID()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al consultar el usuario"})
	}

	//fmt.Println("Usuario recuperado de la base de datos: ", usuario)

	return c.Status(http.StatusOK).JSON(usuario)

}

func (user *Usuario) CrearUsuario(c *fiber.Ctx) error {

	usuario := models.Usuario{}

	err := c.BodyParser(&usuario)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "No se pudieron convertir los datos"})
	}

	err = usuario.Crear()
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: err.Error()})
	}

	return c.JSON(usuario)
}

// endpoint para subir avatar
func (user *Usuario) subirAvatar(c *fiber.Ctx) error {
	// claims := c.Locals("claims").(*models.Claim)

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

	// result := database.Database.Model(&models.Usuario{ID: claims.IdUser}).Update("URLAvatar", fileAvatarName)
	// if result.Error != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al actualizar el nombre de la imagen"})
	// }

	return c.SendStatus(http.StatusOK)
}

// endpoint para subir avatar en el s3 de amazon
func (user *Usuario) subirAvatarAWS(c *fiber.Ctx) error {
	// claims := c.Locals("claims").(*models.Claim)

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

	// result := database.Database.Model(&models.Usuario{ID: claims.IdUser}).Update("URLAvatar", fileAvatarName)
	// if result.Error != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al actualizar el nombre de la imagen"})
	// }

	return c.SendStatus(http.StatusOK)
}

//Agregar un nuevo trabajo al usuario que se pasa por el id
func (u *Usuario) AgregarTrabajo(c *fiber.Ctx) error {
	var trabajo models.Trabajo

	idusuario := c.Params("idusuario")
	ID, err := strconv.ParseInt(idusuario, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	if err = c.BodyParser(&trabajo); err != nil {
		return c.Status(400).SendString("Error al convertir los datos")
	}
	fmt.Println("Trabajo: ", trabajo)
	err = models.Usuario{ID: uint64(ID)}.AgregarTrabajo(trabajo)
	if err != nil {
		return c.Status(400).SendString("Error al registrar el trabajo")
	}

	return c.SendStatus(http.StatusOK)
}

// TODO: Actualizar usuario
func (u *Usuario) Actualizar(c *fiber.Ctx) error {
	usuario := models.Usuario{}
	idusuario := c.Params("id")
	ID, err := strconv.ParseInt(idusuario, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	if err = c.BodyParser(&usuario); err != nil {
		return c.Status(400).SendString("Error al convertir los datos")
	}

	usuario.ID = uint64(ID)

	err = usuario.Actualizar()
	if err != nil {
		return c.Status(400).SendString("Error al actualizar en la DB " + err.Error())
	}

	return c.JSON(usuario)
}

// confirma un usuario, que verifica su usuario correo
func (u *Usuario) ConfirmarCorreo(c *fiber.Ctx) error {
	usuario := models.Usuario{}
	idusuario := c.Params("id")
	ID, err := strconv.ParseInt(idusuario, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	usuario.ID = uint64(ID)

	err = usuario.ConfirmarCorreo()
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	return c.SendStatus(http.StatusOK)
}
