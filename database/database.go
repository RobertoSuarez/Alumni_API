package database

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetInstancia(config *viper.Viper) *gorm.DB {
	pgUser := config.GetString("POSTGRES_USER")
	pgPassword := config.GetString("POSTGRES_PASSWORD")
	pgDB := config.GetString("POSTGRES_DB")
	pgPort := config.GetString("POSTGRES_PORT")
	pgHost := config.GetString("POSTGRES_HOST")
	appAttempts := config.GetInt("APP_ATTEMPTS")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Guayaquil", pgHost, pgUser, pgPassword, pgDB, pgPort)
	fmt.Println("dsn de postgres:", dsn)
	// El servidor intentara conectarse a la base de datos
	// la cantidad de intentos que se establecio en las variables
	// de entorno.
	var (
		db  *gorm.DB
		err error
	)

	for i := 0; i < appAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			// en caso de no conectarse de espera por 3 segundos
			log.Println(err)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}

	return db
}
