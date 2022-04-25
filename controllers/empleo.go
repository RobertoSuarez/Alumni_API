package controllers

import (
	"fmt"
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

	app.Get("/", e.Logs_busquedas, e.ListarEmpleos)
	app.Post("/", e.Crear)
	app.Get("/autocompletado", e.EmpleoAutocompletado)

	app.Post("/guardados-id", ValidarJWT, e.EmpleosGuardados) // Empleos guardados pero solo retorna un slice de ids
	app.Get("/guardados", ValidarJWT, e.ObtenerEmpleosGuardados)
	app.Post("/:id/guardar", ValidarJWT, e.GuardarEmpleoParaUsuario) // el empleo se guardara para el usuario
	app.Delete("/:id/guardar", ValidarJWT, e.EliminarEmpleoGuardado) // Remover el empleo guardado por el usuario

	// Aplicaciónes de empleos
	app.Get("/aplicar", ValidarJWT, e.ObtenerEmpleosAplicados)
	app.Post("/:id/aplicar", ValidarJWT, e.AplicarEmpleo)              // aplica a un empleo, este metodo anteriormente estaba en usuario
	app.Delete("/:id/aplicar", ValidarJWT, e.EliminarAplicacionEmpleo) // Eliminar aplicación de empleo
	app.Post("/:id/aplicar/estado", ValidarJWT, e.EstadoAplicacion)    // Revisa el estado de la aplicación de empleo

	app.Put("/:id", e.Actualizar)
	app.Get("/:id", e.ObtenerEmpleoByID)
	return app
}

// obtener todos los empleos registrador, este endpoint se
// utiliza en el buscador de la app
func (e *Empleo) ListarEmpleos(c *fiber.Ctx) error {

	maps := make(map[string]interface{})

	titulo := c.Query("titulo")
	// if len(titulo) > 0 {
	// 	maps["titulo"] = titulo //[]string{"Adminstrador de empresa", "Desarrollador de software"}
	// }

	ciudad_id := c.Query("ciudad_id")
	if len(ciudad_id) > 0 {
		maps["ciudad_id"] = ciudad_id
	}

	area := c.Query("area_id")
	if len(area) > 0 {
		maps["area_id"] = area
	}

	page, _ := strconv.Atoi(c.Query("page", "0"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size", "0"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	fmt.Println(page, pageSize)
	// cantidad de registros a saltar.
	offset := (page - 1) * pageSize

	empleos, err := models.Empleo{}.ObtenerTodos(offset, pageSize, maps, titulo)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.JSON(empleos)
}

// registra busquedas
func (e *Empleo) Logs_busquedas(c *fiber.Ctx) error {
	log := models.LogBusquedas{
		Titulo:   c.Query("titulo"),
		CiudadID: ParseInt(c.Query("ciudad_id")),
		AreaId:   ParseInt(c.Query("area_id")),
	}
	go log.Guardar()
	return c.Next()
}

func ParseInt(numero string) uint64 {
	ID, err := strconv.ParseInt(numero, 10, 64)
	if err != nil {
		return 0
	}
	return uint64(ID)
}

// Palabras mas buscadas
func (e *Empleo) EmpleoAutocompletado(c *fiber.Ctx) error {
	titulo := c.Query("titulo")
	//fmt.Println(titulo)
	titulos := models.LogBusquedas{}.Autocompletado(titulo)

	return c.JSON(titulos)
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
// Estos empleos son los guardados por el usuario
func (Empleo) ObtenerEmpleosGuardados(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)

	usuario := models.Usuario{ID: claims.IdUser}

	empleos, err := usuario.ObtenerEmpleosGuardados()
	if err != nil {
		return c.Status(400).SendString("Error al consultar los empleos guardados")
	}

	return c.JSON(empleos)
}

// Eliminar empleo guardado por el usuario
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

func (Empleo) AplicarEmpleo(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)
	usuario := models.Usuario{ID: claims.IdUser}

	ID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	err = usuario.AplicarEmpleo(uint64(ID))
	if err != nil {
		return c.Status(400).SendString("Error al aplicar a este empleo")
	}

	return c.Status(http.StatusOK).SendString("Perfecto ya aplicastes para este trabajo")
}

// Listar todos los empleos aplicados
func (Empleo) ObtenerEmpleosAplicados(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)
	usuario := models.Usuario{ID: claims.IdUser}

	empleos, err := usuario.ObtenerEmpleosAplicados()
	if err != nil {
		return c.Status(400).SendString("Error al consultar los empleos aplicados")
	}

	return c.JSON(empleos)
}

// Eliminar aplicación de empleo
func (Empleo) EliminarAplicacionEmpleo(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*models.Claim)
	usuario := models.Usuario{ID: claims.IdUser}

	ID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	err = usuario.EliminarEmpleoAplicado(uint64(ID))
	if err != nil {
		return c.Status(400).SendString("Empleo removido correctamente")
	}

	return c.SendStatus(http.StatusOK)
}

// Este endpoint revisara si el usuario ha aplicado al empleo,
// si en la base de dato no existe ningun registro de estos, retorna un false, como no aplicado
// en caso de que si exista un registro, se retorna un true, como si aplicado
func (Empleo) EstadoAplicacion(c *fiber.Ctx) error {
	estado := struct {
		Estado bool `json:"estado"`
	}{}
	claims := c.Locals("claims").(*models.Claim)
	usuario := models.Usuario{ID: claims.IdUser}

	ID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error en el ID")
	}

	estado.Estado = usuario.EstadoAplicacion(uint64(ID))

	return c.JSON(estado)
}
