package models

type OfertaLaboral struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	UsuarioID uint    `json:"usuarioID"`
	Usuario   Usuario `gorm:"foreignKey:UsuarioID"`

	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
}
