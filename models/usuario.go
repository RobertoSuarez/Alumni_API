package models

import "errors"

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
	"RoleCuenta",
	"EmailConfirmado",
	"IsStaff",
	"StaffRole",
	"RoleCuenta",
}

type Usuario struct {
	ID                   uint   `json:"id" gorm:"primary_key"`
	IdentificacionTipo   string `json:"identificacionTipo"`
	NumeroIdentificacion string `json:"numeroIdentificacion"`
	Nombres              string `json:"nombres"`
	Apellidos            string `json:"apellidos"`
	Email                string `json:"email"`
	Password             string `json:"password,omitempty"`
	Nacimiento           string `json:"nacimiento"`
	Whatsapp             string `json:"whatsapp"`
	// TipoUsuarioID        uint        `json:"tipoUsuarioID"`
	// TipoUsuario          TipoUsuario `json:"tipoUsuario" gorm:"foreignKey:TipoUsuarioID"`

	EmailConfirmado bool `json:"emailConfirmado,omitempty"`

	IsStaff   bool   `json:"isStaff"`   // si es del staff
	StaffRole string `json:"staffRole"` // Administrador, moderador u otra cosa

	RoleCuenta string `json:"roleCuenta"` // si es alumno, alumni o usuarionormal

	// // agregar los tipos de usuarios
	// AdminID uint   `json:"-"`
	// Admin   *Admin `json:"admin,omitempty" gorm:"foreignKey:AdminID"`

	// AlumniID uint    `json:"-"`
	// Alumni   *Alumni `json:"alumni,omitempty" gorm:"foreignKey:AlumniID"`

	// datos de gorm

	OfertasLaborales []OfertaLaboral `json:"ofertaLaboral,omitempty" gorm:"foreignKey:UsuarioID"`

	Educacion []Educacion `json:"educacion,omitempty" gorm:"foreignKey:UsuarioID"` // este va hacer el historial academico
}

type TipoUsuario struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Tipo string `json:"tipo"`
}

type Estado struct {
	Usando bool
}

type ListTipoUsuarios []TipoUsuario

func (listTipos *ListTipoUsuarios) GetID(tipo string) (uint, error) {
	for _, v := range *listTipos {
		if v.Tipo == tipo {
			return v.ID, nil
		}
	}

	return 0, errors.New("no existe ese tipo de usuario")
}
