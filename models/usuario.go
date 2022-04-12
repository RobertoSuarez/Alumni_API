package models

import (
	"errors"
	"fmt"
	"time"
)

var UsuarioCamposDB = []string{
	"ID",
	"IdentificacionTipo",
	"NumeroIdentificacion",
	"Nombres",
	"Apellidos",
	"Email",
	//"Password",
	"Nacimiento",
	"Whatsapp",
	"URLAvatar",
	"RoleCuenta",
	"EmailConfirmado",
	"IsStaff",
	"StaffRole",
	"RoleCuenta",
}

const datos = "id,created_at,updated_at,identificacion_tipo,numero_identificacion,nombres,apellidos,email,nacimiento,phone,avatar,descripcion,is_super,email_confirmado,genero,fecha_graduacion,nivel_academico,es_discapacitado"

type Usuario struct {
	ID                   uint64 `json:"id" gorm:"primary_key"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	IdentificacionTipo   string    `json:"identificacionTipo" gorm:"size:100"`
	NumeroIdentificacion string    `json:"numeroIdentificacion" gorm:"size:100"`
	Nombres              string    `json:"nombres" gorm:"size:200"`
	Apellidos            string    `json:"apellidos" gorm:"size:200"`
	Email                string    `json:"email" gorm:"not null;unique"`
	Password             string    `json:"password" gorm:"size:200;not null"`
	Nacimiento           time.Time `json:"nacimiento"`
	Phone                string    `json:"phone"`
	Avatar               string    `json:"avatar"`
	Descripcion          string    `json:"descripcion"`
	IsSuper              bool      `json:"is_super_user"`
	EmailConfirmado      bool      `json:"emailConfirmado"`
	Genero               string    `json:"genero" gorm:"size:75"`
	FechaGraduacion      time.Time `json:"fechaGraduacion"`
	NivelAcademico       string    `json:"nivelAcademico"`
	EsDiscapacitado      bool      `json:"esDiscapacitado"`
	Grupos               []Grupo   `json:"grupos,omitempty" gorm:"many2many:usuario_grupos;"`
	Trabajos             []Trabajo `json:"trabajos,omitempty" gorm:"foreignKey:UsuarioID"`
	EmpresasPropias      []Empresa `json:"empresasPropias,omitempty" gorm:"foreignKey:UsuarioCreadorID"`
	EmpresasAsociadas    []Empresa `json:"empresaAsociadas,omitempty" gorm:"many2many:usuario_empresas_asociadas;"`
	EmpleosGuardados     []Empleo  `json:"empleosGuardados" gorm:"many2many:empleos_guardado"`
	EmpleosAplicados     []Empleo  `json:"empleosAplicados" gorm:"many2many:empleos_aplicados"`
}

func (Usuario) TableName() string {
	return "usuario"
}

// LoginUsuario revisar en la db que el usuario y contreseña existan.
func (Usuario) LoginUsuario(login Login) (Usuario, error) {
	usuario := Usuario{}

	// Busca el usuario con ese email
	result := DB.Where("email = ?", login.Email).First(&usuario)

	// controlamos el error
	if result.Error != nil {
		return usuario, errors.New("no existe el usuario")
	}

	if usuario.Password != login.Password {
		return usuario, errors.New("la contraseña es incorrecta")
	}

	// Traemos todos los grupos del usuario
	DB.Model(&usuario).Association("Grupos").Find(&usuario.Grupos)

	return usuario, nil
	//Database.Preload("Admin").Preload("Alumni").Preload("TipoUsuario").First(&usuario)
}

// Obtener todos los usuario registrado en la base de datos
func (Usuario) GetAll() (usuarios []Usuario, err error) {
	usuarios = []Usuario{}
	result := DB.Select(datos).Find(&usuarios)
	if result.Error != nil {
		fmt.Println(result.Error)
		return usuarios, result.Error
	}

	return usuarios, nil
}

// GetUsuarioByID retorna el usuario mediante el ID, el ID se
// debe pasar con el objeto que lo invoca.
func (u Usuario) GetUsuarioByID() (usuario Usuario, err error) {
	err = DB.Omit("Password").First(&usuario, u.ID).Error
	if err != nil {
		return usuario, err
	}

	DB.Model(&usuario).Association("Grupos").Find(&usuario.Grupos)
	DB.Model(&usuario).Association("Trabajos").Find(&usuario.Trabajos)
	DB.Model(&usuario).Association("EmpresasPropias").Find(&usuario.EmpresasPropias)
	return usuario, nil
}

// Crear metodo para insertar un usuario en la base de datos
func (u *Usuario) Crear() error {
	tx := DB.Begin()

	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Agregar grupo a usuario
func (u Usuario) AgregarGrupo(g Grupo) error {
	tx := DB.Begin()

	err := tx.Model(&u).Association("Grupos").Append(&g)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//Agregar trabajo agregara un nuevo registro de trabajo
// que solo este mismo usuario prodra administrar
func (u Usuario) AgregarTrabajo(t Trabajo) error {
	tx := DB.Begin()
	fmt.Println("Trabajo en models: ", t)
	err := tx.Model(&u).Association("Trabajos").Append(&t)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Elimina un usuario en todas sus instancias
func (u Usuario) Eliminar(ids []uint64) error {
	tx := DB.Begin()

	if err := tx.Exec("DELETE FROM usuario WHERE id IN (?)", ids).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Exec("DELETE FROM usuario_grupos WHERE usuario_id IN (?)", ids).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (u *Usuario) Actualizar() error {
	tx := DB.Begin()
	//campos := "identificacion_tipo,numero_identificacion,nombres,apellidos,nacimiento,phone,descripcion,genero,fecha_graduacion,nivel_academico,es_discapacitado"
	// Campos que se actualizaran en la tabla
	// se actualizaran los datos, no relevantes
	err := tx.Model(&u).Updates(Usuario{
		IdentificacionTipo:   u.IdentificacionTipo,
		NumeroIdentificacion: u.NumeroIdentificacion,
		Nombres:              u.Nombres,
		Apellidos:            u.Apellidos,
		Nacimiento:           u.Nacimiento,
		Phone:                u.Phone,
		Genero:               u.Genero,
		FechaGraduacion:      u.FechaGraduacion,
		NivelAcademico:       u.NivelAcademico,
		EsDiscapacitado:      u.EsDiscapacitado,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.First(&u, u.ID)

	tx.Commit()
	return nil
}

// Actualizar la descripción
func (u Usuario) ActualizarDescripcion() error {
	tx := DB.Begin()

	result := tx.Model(&u).Updates(Usuario{
		Descripcion: u.Descripcion,
	})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}

// Verificar o confirmar el correo de la cuenta
func (u Usuario) ConfirmarCorreo() error {

	if u.ID < 1 {
		return errors.New("falta el id del usuario")
	}
	tx := DB.Begin()

	result := tx.Model(&u).Update("email_confirmado", true)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected < 1 {
		tx.Rollback()
		return errors.New("no existe usuario con ese id")
	}

	tx.Commit()
	return nil
}

// Guardar empleo
// Usuaruio guarda un empleo
func (u Usuario) GuardarEmpleo(id_empleo uint64) (Empleo, error) {
	empleo := Empleo{ID: id_empleo}

	tx := DB.Begin()

	err := tx.Model(&u).Association("EmpleosGuardados").Append(&empleo)
	if err != nil {
		tx.Rollback()
		return empleo, err
	}

	tx.Model(&empleo).Preload("Area").Preload("Subarea").First(&empleo)

	tx.Commit()
	return empleo, nil
}

// Listar todos los empleos guardados
func (u Usuario) ObtenerEmpleosGuardados() ([]Empleo, error) {

	empleos := []Empleo{}

	err := DB.Model(&u).Preload("Area").Preload("Subarea").Association("EmpleosGuardados").Find(&empleos)
	if err != nil {
		return empleos, err
	}

	return empleos, nil
}

// Eliminar empleo guardado
func (u Usuario) EliminarEmpleoGuardado(id_empleo uint64) error {
	tx := DB.Begin()

	err := tx.Model(&u).Association("EmpleosGuardados").Delete(&Empleo{ID: id_empleo})
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Obtener los empleos guardados, pero de la lista de empleos que me envia el cliente (request)
func (u Usuario) ObtenerEmpleosGuardadosIDVerificar(ids []uint64) ([]uint64, error) {
	empleosID := []uint64{}

	err := DB.Model(&u).Where("empleo_id IN ?", ids).Select("empleo_id").Association("EmpleosGuardados").Find(&empleosID)
	if err != nil {
		return empleosID, err
	}

	return empleosID, nil
}

// Este usuario aplicara a un empleo
func (u Usuario) AplicarEmpleo(id_empleo uint64) error {
	tx := DB.Begin()

	err := tx.Model(&u).Association("EmpleosAplicados").Append(&Empleo{ID: id_empleo})
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Listar todos los empleos a los cuales a aplicado
func (u Usuario) ObtenerEmpleosAplicados() ([]Empleo, error) {

	empleos := []Empleo{}

	err := DB.Model(&u).Association("EmpleosAplicados").Find(&empleos)
	if err != nil {
		return empleos, err
	}

	return empleos, nil
}

// Eliminar empleos aplicados
func (u Usuario) EliminarEmpleoAplicado(id_empleo uint64) error {
	tx := DB.Begin()

	err := tx.Model(&u).Association("EmpleosAplicados").Delete(&Empleo{ID: id_empleo})
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (u Usuario) EstadoAplicacion(id_empleo uint64) bool {

	count := DB.Model(&u).Where("empleo_id = ?", id_empleo).Association("EmpleosAplicados").Count()

	fmt.Println("Cantidad de Empleos asociados pero con el where: ", count)

	return count == 1

}
