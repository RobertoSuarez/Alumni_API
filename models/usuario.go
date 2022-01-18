package models

type Usuario struct {
	ID                   uint        `json:"id" gorm:"primary_key"`
	IdentificacionTipo   string      `json:"identificacionTipo"`
	NumeroIdentificacion string      `json:"numeroIdentificacion"`
	Nombres              string      `json:"nombres"`
	Apellidos            string      `json:"apellidos"`
	Email                string      `json:"email"`
	Password             string      `json:"password,omitempty"`
	Nacimiento           string      `json:"nacimiento"`
	Whatsapp             string      `json:"whatsapp"`
	TipoUsuarioID        uint        `json:"tipoUsuarioID"`
	TipoUsuario          TipoUsuario `json:"tipoUsuario"`

	// agregar los tipos de usuarios
	AdminID uint   `json:"-"`
	Admin   *Admin `json:"admin,omitempty" gorm:"foreignKey:AdminID"`

	AlumniID uint    `json:"-"`
	Alumni   *Alumni `json:"alumni,omitempty" gorm:"foreignKey:AlumniID"`
}

type TipoUsuario struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Tipo string `json:"tipo"`
}

type PerfilCompleto struct {
	Usuario Usuario `json:"usuario"`
	Admin   Admin   `json:"admin,omitempty"`
	Alumni  Alumni  `json:"alumni,omitempty"`
}

type Estado struct {
	Usando bool
}
