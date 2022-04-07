package models

import "time"

type Area struct {
	ID          uint64 `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Titulo      string `json:"titulo" gorm:"size:200"`
	Descripcion string `json:"descripcion" gorm:"size:900"`
}

func (Area) TableName() string {
	return "area"
}

// Listar las areas
func (Area) ObtenerAreas() ([]Area, error) {
	areas := []Area{}

	result := DB.Find(&areas)
	if result.Error != nil {
		return areas, result.Error
	}

	return areas, nil
}
