package models

type RoleUsuario struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Titulo      string `json:"titulo" gorm:"size:100"`
	Descripcion string `json:"descripcion" gorm:"size:300"`
}

func (RoleUsuario) TableName() string {
	return "RoleUsuario"
}
