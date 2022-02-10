package models

type Alumni struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Carrera string `json:"carrera"`
	Estado  `json:"-"`
}

func (alumni *Alumni) SetNil() *Alumni {
	alumni.Estado = Estado{Usando: false}
	alumni.Carrera = "--- nil"
	return alumni
}
