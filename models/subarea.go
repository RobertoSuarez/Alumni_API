package models

import "time"

type Subarea struct {
	ID          uint64 `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Titulo      string `json:"titulo" gorm:"size:200"`
	Descripcion string `json:"descripcion" gorm:"size:900"`
	AreaID      uint64
	Area        Area `json:"area" gorm:"foreignKey:AreaID"`
}

func (Subarea) TableName() string {
	return "subarea"
}
