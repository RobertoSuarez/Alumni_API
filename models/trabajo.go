package models

import "time"

type Trabajo struct {
	ID            uint64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Area          string `json:"area" gorm:"size:200"`
	Cargo         string `json:"cargo" gorm:"size:200"`
	TrabajoActual bool   `json:"trabajoActual"`
	UsuarioID     uint64
}

func (Trabajo) TableName() string {
	return "trabajo"
}
