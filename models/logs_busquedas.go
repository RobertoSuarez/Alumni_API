package models

type LogBusquedas struct {
	ID       uint64 `json:"id" gorm:"primary_key"`
	Titulo   string `json:"titulo"`
	CiudadID uint64 `json:"ciudad_id"`
	AreaId   uint64 `json:"area_id"`
}

func (LogBusquedas) TableName() string {
	return "logs_busquedas"
}

func (l LogBusquedas) Guardar() {
	DB.Create(&l)
}

func (l LogBusquedas) Autocompletado(titulo string) []string {
	busquedas := []string{}
	DB.Model(&l).
		Limit(5).
		Where("titulo LIKE ?", "%"+titulo+"%").
		Pluck("titulo", &busquedas)
	return busquedas
}
