package models

type Alumni struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Carrera string `json:"carrera"`
	Estado  `json:"-"`
}
