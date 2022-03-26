package models

type Grupo struct {
	ID       uint64    `json:"id" gorm:"primary_key"`
	Name     string    `json:"name" gorm:"not null;unique"`
	Usuarios []Usuario `json:"usuario" gorm:"many2many:usuario_grupos;"`
}

func (Grupo) TableName() string {
	return "grupo"
}
