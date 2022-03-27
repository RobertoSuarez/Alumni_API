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
