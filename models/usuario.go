package models

import (
	"errors"
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
	ID                   uint      `json:"id" gorm:"primary_key"`
	IdentificacionTipo   string    `json:"identificacionTipo" gorm:"size:100"`
	NumeroIdentificacion string    `json:"numeroIdentificacion" gorm:"size:100"`
	Nombres              string    `json:"nombres" gorm:"size:200"`
	Apellidos            string    `json:"apellidos" gorm:"size:200"`
	Email                string    `json:"email" gorm:"size:200;unique;not null"`
	Password             string    `json:"password,omitempty" gorm:"size:200"`
	Nacimiento           time.Time `json:"nacimiento"`
	Whatsapp             string    `json:"whatsapp" gorm:"size:200"`
	URLAvatar            string    `json:"urlAvatar"`
	Descripcion          string    `json:"descripcion"`
	// TipoUsuarioID        uint        `json:"tipoUsuarioID"`
	// TipoUsuario          TipoUsuario `json:"tipoUsuario" gorm:"foreignKey:TipoUsuarioID"`

	EmailConfirmado bool `json:"emailConfirmado,omitempty"`

	IsStaff   bool   `json:"isStaff"`                   // si es del staff
	StaffRole string `json:"staffRole" gorm:"size:200"` // Administrador, moderador u otra cosa

	RoleCuenta string `json:"roleCuenta" gorm:"size:200"` // si es alumno, alumni o usuarionormal

	// // agregar los tipos de usuarios
	// AdminID uint   `json:"-"`
	// Admin   *Admin `json:"admin,omitempty" gorm:"foreignKey:AdminID"`

	// AlumniID uint    `json:"-"`
	// Alumni   *Alumni `json:"alumni,omitempty" gorm:"foreignKey:AlumniID"`

	// datos de gorm

	OfertasLaborales []Empleo `json:"ofertaLaboral,omitempty" gorm:"foreignKey:UsuarioID"`

	Educacion []Educacion `json:"educacion,omitempty" gorm:"foreignKey:UsuarioID"` // este va hacer el historial academico
}

func (Usuario) TableName() string {
	return "usuario"
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
