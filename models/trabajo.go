package models

import (
	"errors"
	"time"
)

type Trabajo struct {
	ID            uint64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Area          string `json:"area" gorm:"size:200"`
	Cargo         string `json:"cargo" gorm:"size:200"`
	TrabajoActual bool   `json:"trabajoActual"`
	UsuarioID     uint64
}

func (Trabajo) TableName() string {
	return "trabajo"
}

// Todo: listar todos los trabajo por que esten filtrado por el id usuario
func (Trabajo) ObtenerTrabajosUsuario(id uint64) (trabajos []Trabajo, err error) {
	result := DB.Where("usuario_id = ?", id).Find(&trabajos)
	if result.Error != nil {
		return trabajos, result.Error
	}

	return trabajos, nil
}

// TODO: Actualizar trabajo
func (t *Trabajo) Actualizar() error {

	tx := DB.Begin()
	result := tx.Model(&t).Updates(Trabajo{
		Area:          t.Area,
		Cargo:         t.Cargo,
		TrabajoActual: t.TrabajoActual,
	})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected < 1 {
		tx.Rollback()
		return errors.New("no existe registro")
	}

	tx.Commit()

	return nil
}

// TDDO: Eliminar trabajo
func (t *Trabajo) Eliminar() error {

	result := DB.Delete(&t)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected < 1 {
		return errors.New("no existe el registro")
	}

	return nil
}
