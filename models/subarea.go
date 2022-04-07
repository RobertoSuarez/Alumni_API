package models

import "time"

type Subarea struct {
	ID          uint64 `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Titulo      string `json:"titulo" gorm:"size:200"`
	Descripcion string `json:"descripcion" gorm:"size:900"`
	AreaID      uint64 `json:"areaid,omitempty"`
	Area        *Area  `json:"area,omitempty" gorm:"foreignKey:AreaID"`
}

func (Subarea) TableName() string {
	return "subarea"
}

// Listar las subareas en base al area
func (Subarea) ObtenerSubareas(areaID uint64) ([]Subarea, error) {
	areas := []Subarea{}
	result := DB.Where("area_id = ?", areaID).Find(&areas)
	if result.Error != nil {
		return areas, result.Error
	}

	return areas, nil
}
