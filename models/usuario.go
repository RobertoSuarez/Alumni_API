package models

import "gorm.io/gorm"

type Usuario struct {
	gorm.Model            `json:"-"`
	Identificacion_tipo   string `json:"identificacionTipo"`
	Numero_identificacion string `json:"numeroIdentificacion"`
	Nombres               string `json:"nombres"`
	Apellidos             string `json:"apellidos"`
	Email                 string `json:"email"`
	Password              string `json:"password,omitempty"`
	Nacimiento            string `json:"nacimiento"`
	Whatsapp              string `json:"whatsapp"`
	TipoUsuario           string `json:"tipoUsuario"`
}
