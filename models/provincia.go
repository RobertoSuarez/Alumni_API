package models

type Provincia struct {
	ID       uint64
	Nombre   string
	Ciudades []Ciudad `json:"ciudades" gorm:"foreignKey:ProvinciaID"`
}

func (Provincia) TableName() string {
	return "provincia"
}

func (Provincia) ObtenerTodas() ([]Provincia, error) {
	provincias := []Provincia{}
	result := DB.Find(&provincias)
	if result.Error != nil {
		return provincias, result.Error
	}

	return provincias, nil
}

// Obtener las ciudades por el id de la provincia
func (pro Provincia) ObtenerCiudades() ([]Ciudad, error) {
	ciudades := []Ciudad{}

	err := DB.Model(&pro).Association("Ciudades").Find(&ciudades)
	if err != nil {
		return ciudades, err
	}

	return ciudades, nil

}
