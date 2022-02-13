package models

import "time"

type Educacion struct {
	ID                    uint      `json:"id" gorm:"primary_key"`
	UsuarioID             uint      `json:"usuarioID"` // foreignKey
	InstituacionEducativa string    `json:"instituacionEducativa"`
	Titulo                string    `json:"titulo"`
	DiciplinaAcademica    string    `json:"diciplinaAcademica"`
	FechaInicio           time.Time `json:"fechaInicio"`
	FechaFin              time.Time `json:"fechaFin"`
	Nota                  string    `json:"nota"`
	Descripcion           string    `json:"descripcion"`
}
