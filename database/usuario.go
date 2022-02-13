package database

import (
	"errors"

	"github.com/RobertoSuarez/apialumni/models"
)

// LoginUsuario revisar en la db que el usuario y contreseña existan.
func LoginUsuario(email, password string) (*models.Usuario, error) {
	usuario := models.Usuario{}

	// Busca en la base de datos
	result := Database.
		Where("email = ? AND password = ?", email, password).First(&usuario)

	// controlamos el error
	if result.Error != nil {
		return nil, errors.New("no existe ese registro")
	}

	return &usuario, nil
	//Database.Preload("Admin").Preload("Alumni").Preload("TipoUsuario").First(&usuario)
}
