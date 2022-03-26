package database

import (
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetInstancia(config *viper.Viper) *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.GetString("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
