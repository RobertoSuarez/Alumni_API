package models

import "time"

type OfertaLaboral struct {
	ID        uint     `json:"id" gorm:"primary_key"`
	UsuarioID uint     `json:"usuarioID"` // foreignKey
	Usuario   *Usuario `json:"usuario,omitempty"`

	Fecha   time.Time `json:"fecha"`
	Empresa string    `json:"empresa"` // TODO: sdfjk

	Titulo                   string `json:"titulo"`
	Descripcion              string `json:"descripcion"`
	Profesion                string `json:"profesion"`
	Puesto                   string `json:"puesto"`
	TipoEmplo                string `json:"tipoEmpleo"` //Modalidad de trabajo
	Area                     string `json:"area"`       // Categoria
	Subarea                  string `json:"subarea"`
	Sueldo                   string `json:"sueldo"`
	TiempoExperiencia        string `json:"tiempoExperiencia"` // Los años de experiencia
	Jornada                  string `json:"jornada"`
	TipoContrato             string `json:"tipoContrato"`
	ConocimientosAdicionales string `json:"conocimientosAdicionales"`
	Ciudad                   string `json:"ciudad"`
	PostulanteDiscapacidad   bool   `json:"postulanteDiscapacidad"` // si el trabajo es para personas con capacidades limitadas.
}
