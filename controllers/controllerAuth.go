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

	app.Get("/users", ValidarJWT, gruposPermitios([]string{"estudiante"}), func(c *fiber.Ctx) error {
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
	login := models.Login{}

	err := c.BodyParser(&login)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al convertir los datos"})
	}

	usuario, err := models.Usuario{}.LoginUsuario(login)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	token, err := GenerarJWT(usuario)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorAPI{Mensaje: "Error al construir el token"})
	}

	usuario.Password = ""

	respuestaLogin := models.RespuestaLogin{
		Token:   token,
		Usuario: usuario,
	}

	return c.JSON(respuestaLogin)
}
