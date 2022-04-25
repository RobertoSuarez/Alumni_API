package models

type Ciudad struct {
	ID          uint64
	Nombre      string
	ProvinciaID uint64
	Provincia   Provincia `json:"-" gorm:"foreignKey:ProvinciaID"`
}

func (Ciudad) TableName() string {
	return "ciudad"
}

// Listar las ciudades
func (Ciudad) ObtenerTodas() ([]Ciudad, error) {
	ciudades := []Ciudad{}
	result := DB.Find(&ciudades)
	if result.Error != nil {
		return ciudades, result.Error
	}

	return ciudades, nil
}
