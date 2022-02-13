package models

type OfertaLaboral struct {
	ID          uint    `json:"id" gorm:"primary_key"`
	UsuarioID   uint    `json:"usuarioID"` // foreignKey
	Usuario     Usuario `json:"usuario,omitempty"`
	Titulo      string  `json:"titulo"`
	Descripcion string  `json:"descripcion"`
}
