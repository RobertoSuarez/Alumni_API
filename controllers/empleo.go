package controllers

import (
	"net/http"
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

	app.Post("/guardados-id", ValidarJWT, e.EmpleosGuardados) // Empleos guardados pero solo retorna un slice de ids
	app.Get("/guardados", ValidarJWT, e.ObtenerEmpleosGuardados)
	app.Post("/:id/guardar", ValidarJWT, e.GuardarEmpleoParaUsuario) // el empleo se guardara para el usuario
	app.Delete("/:id/guardar", ValidarJWT, e.EliminarEmpleoGuardado) // Remover el empleo guardado por el usuario

	app.Put("/:id", e.Actualizar)
	app.Get("/:id", e.ObtenerEmpleoByID)
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

// Trabajos guardados verificación, esto debe recibir todos los trabajos id,
// y retornar los que si estan guardados, se debe tener claro que este endpoint
// no guarda ningun trabajo, solo verifica si el usuario ya lo tiene guardado
func (Empleo) EmpleosGuardados(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)
	usuario := models.Usuario{ID: claims.IdUser}

	ids := []uint64{}
	err := c.BodyParser(&ids)
	if err != nil {
		return c.Status(400).SendString("No se pudo convertir los ids")
	}

	IDs, err := usuario.ObtenerEmpleosGuardadosIDVerificar(ids)
	if err != nil {
		return c.Status(400).SendString("Error al aplicar a este empleo")
	}

	return c.JSON(IDs)
}

// Obtener empleo por el id
func (Empleo) ObtenerEmpleoByID(c *fiber.Ctx) error {
	idempleo := c.Params("id")
	ID, err := strconv.ParseInt(idempleo, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}
	empleo := models.Empleo{ID: uint64(ID)}

	err = empleo.ObtenerEmpleoByID()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	return c.JSON(empleo)
}

// GuardarEmpleoParaUsuario Relaciona el usuario del token con
// el id del empleo que se pasa
func (Empleo) GuardarEmpleoParaUsuario(c *fiber.Ctx) error {
	// Obtenemos los privilegio y el id como parametro
	claims := c.Locals("claims").(*models.Claim)
	usuario := models.Usuario{ID: claims.IdUser}
	idempleo := c.Params("id")
	ID, err := strconv.ParseInt(idempleo, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	// Ingresamos la relación
	empleo, err := usuario.GuardarEmpleo(uint64(ID))
	if err != nil {
		return c.Status(400).SendString("No se puddo guardar")
	}

	return c.JSON(empleo)
}

// ObtenerEmpleosGuardados obtien los emplos guardados pero todo el objeto
func (Empleo) ObtenerEmpleosGuardados(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)

	usuario := models.Usuario{ID: claims.IdUser}

	empleos, err := usuario.ObtenerEmpleosGuardados()
	if err != nil {
		return c.Status(400).SendString("Error al consultar los empleos guardados")
	}

	return c.JSON(empleos)
}

func (Empleo) EliminarEmpleoGuardado(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)

	ID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	usuario := models.Usuario{ID: claims.IdUser}

	err = usuario.EliminarEmpleoGuardado(uint64(ID))
	if err != nil {
		return c.Status(400).SendString("No se puedo eliminar")
	}

	return c.SendStatus(200)
}
