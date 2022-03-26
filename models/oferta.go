package models

import "time"

type Empleo struct {
	ID        uint     `json:"id" gorm:"primary_key"`
	UsuarioID uint64   `json:"usuarioID"` // foreignKey
	Usuario   *Usuario `json:"usuario,omitempty"`

	Fecha   time.Time `json:"fecha"`
	Empresa string    `json:"empresa" gorm:"size:200"` // TODO: sdfjk

	Titulo                   string `json:"titulo" gorm:"size:200"`
	Descripcion              string `json:"descripcion"`
	Profesion                string `json:"profesion" gorm:"size:200"`
	Puesto                   string `json:"puesto" gorm:"size:200"`
	TipoEmplo                string `json:"tipoEmpleo" gorm:"size:200"` //Modalidad de trabajo
	Area                     string `json:"area" gorm:"size:200"`       // Categoria
	Subarea                  string `json:"subarea" gorm:"size:200"`
	Sueldo                   string `json:"sueldo" gorm:"size:200"`
	TiempoExperiencia        string `json:"tiempoExperiencia" gorm:"size:200"` // Los años de experiencia
	Jornada                  string `json:"jornada" gorm:"size:200"`
	TipoContrato             string `json:"tipoContrato" gorm:"size:200"`
	ConocimientosAdicionales string `json:"conocimientosAdicionales" gorm:"size:200"`
	Ciudad                   string `json:"ciudad" gorm:"size:200"`
	PostulanteDiscapacidad   bool   `json:"postulanteDiscapacidad"` // si el trabajo es para personas con capacidades limitadas.
}

func (Empleo) TableName() string {
	return "empleo"
}

// QueryEmpleo se utilizara para realizar consultas
// En los empleos con una mayor exactitud.
type QueryEmpleo struct {
	Areas     []string `json:"area" query:"areas"`
	Ciudades  []string `json:"ciudades" query:"ciudades"`
	Busquedad string   `json:"busquedad" query:"busquedad"`
}
