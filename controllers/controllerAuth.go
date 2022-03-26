package controllers

import (
	"net/http"

	"github.com/RobertoSuarez/apialumni/models"
	"github.com/gofiber/fiber/v2"
)

type ControllerAuth struct{}

func NewControllerAuth() *ControllerAuth {
	return &ControllerAuth{}
}

// Implementa la interface de <ConfigMicroServicio>
func (c *ControllerAuth) ConfigPath(app *fiber.App) *fiber.App {

	app.Post("/login", c.LoginHandle)

	app.Get("/users", ValidarJWT, func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*models.Claim)
		_ = claims
		//fmt.Println(claims)
		users := []models.Usuario{}

		//database.Database.Find(&users)

		return c.JSON(users)
	})

	return app
}

func (auth *ControllerAuth) GetUsuarios(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)
	_ = claims
	//fmt.Println(claims)
	users := []models.Usuario{}

	//database.Database.Find(&users)

	return c.JSON(users)
}

func (auth *ControllerAuth) LoginHandle(c *fiber.Ctx) error {
	var login models.Login

	err := c.BodyParser(&login)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al convertir los datos"})
	}

	//result := database.Database.Where("email = ? AND password = ?", login.Username, login.Password).First(&usuario)

	// usuario, err := database.LoginUsuario(login.Email, login.Password)

	// if err != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: err.Error()})
	// }

	// token, err := GenerarJWT(usuario)
	// if err != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al construir el token"})
	// }
	// usuario.Password = ""

	// respuestaLogin := models.RespuestaLogin{
	// 	Token:   token,
	// 	Usuario: usuario,
	// }

	//ClearTiposUsuarios(usuario)

	return c.Status(http.StatusOK).JSON(nil)
}
