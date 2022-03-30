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

type Usuario struct {
	ID                   uint64 `json:"id" gorm:"primary_key"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	IdentificacionTipo   string    `json:"identificacionTipo" gorm:"size:100"`
	NumeroIdentificacion string    `json:"numeroIdentificacion" gorm:"size:100"`
	Nombres              string    `json:"nombres" gorm:"size:200"`
	Apellidos            string    `json:"apellidos" gorm:"size:200"`
	Email                string    `json:"email" gorm:"not null;unique"`
	Password             string    `json:"password,omitempty" gorm:"size:200;not null"`
	Nacimiento           time.Time `json:"nacimiento"`
	Phone                string    `json:"phone"`
	Avatar               string    `json:"avatar"`
	Descripcion          string    `json:"descripcion"`
	IsSuper              bool      `json:"is_super_user"`
	EmailConfirmado      bool      `json:"emailConfirmado,omitempty"`
	Genero               string    `json:"genero" gorm:"size:75"`
	FechaGraduacion      time.Time `json:"fechaGraduacion"`
	NivelAcademico       string    `json:"nivelAcademico"`
	EsDiscapacitado      bool      `json:"esDiscapacitado"`
	Grupos               []Grupo   `json:"grupos,omitempty" gorm:"many2many:usuario_grupos;"`
	Trabajos             []Trabajo `json:"trabajos,omitempty" gorm:"foreignKey:UsuarioID"`
	EmpresasPropias      []Empresa `json:"empresasPropias,omitempty" gorm:"foreignKey:UsuarioCreadorID"`
	EmpresasAsociadas    []Empresa `json:"empresaAsociadas,omitempty" gorm:"many2many:usuario_empresas_asociadas;"`
}

func (Usuario) TableName() string {
	return "usuario"
}

// LoginUsuario revisar en la db que el usuario y contreseña existan.
func (Usuario) LoginUsuario(email, password string) (*Usuario, error) {
	usuario := Usuario{}

	// Busca en la base de datos
	result := DB.
		Where("email = ? AND password = ?", email, password).First(&usuario)

	// controlamos el error
	if result.Error != nil {
		return nil, errors.New("no existe ese registro")
	}

	return &usuario, nil
	//Database.Preload("Admin").Preload("Alumni").Preload("TipoUsuario").First(&usuario)
}

// Obtener todos los usuario registrado en la base de datos
func (Usuario) GetAll() (usuarios []Usuario, err error) {
	usuarios = []Usuario{}

	result := DB.Find(&usuarios)
	if result.Error != nil {
		fmt.Println(result.Error)
		return usuarios, result.Error
	}

	return usuarios, nil
}

// GetUsuarioByID retorna el usuario mediante el ID, el ID se
// debe pasar con el objeto que lo invoca.
func (u Usuario) GetUsuarioByID() (usuario Usuario, err error) {
	err = DB.First(&usuario, u.ID).Error
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
	err := tx.Model(&u).Updates(Usuario{
		IdentificacionTipo:   u.IdentificacionTipo,
		NumeroIdentificacion: u.NumeroIdentificacion,
		Nombres:              u.Nombres,
		Apellidos:            u.Apellidos,
		Nacimiento:           u.Nacimiento,
		Phone:                u.Phone,
		Descripcion:          u.Descripcion,
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
