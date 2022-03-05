package database

import (
	"fmt"

	"github.com/RobertoSuarez/apialumni/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDB() {
	viper.AutomaticEnv()
	fmt.Println("url postgresql: ", viper.GetString("DATABASE_URL"))

	db, err := gorm.Open(postgres.Open(viper.GetString("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Database = db

	// Migraciones
	db.AutoMigrate(&models.Usuario{})
	db.AutoMigrate(&models.TipoUsuario{})
	// db.AutoMigrate(&models.Admin{})
	// db.AutoMigrate(&models.Alumni{})
	db.AutoMigrate(&models.Empleo{})
	db.AutoMigrate(&models.Educacion{})

	// tipos de usuarios
	// tipos := []models.TipoUsuario{}
	// result := db.Find(&tipos)
	// if result.Error != nil {
	// 	fmt.Println(result.Error)
	// }
	// if len(tipos) == 0 {
	// 	Database.Create(&models.TipoUsuario{Tipo: "admin"})
	// 	Database.Create(&models.TipoUsuario{Tipo: "alumni"})
	// }

	// db.Create(&models.Usuario{
	// 	Email:    "electrosonix12@gmail.com",
	// 	Password: "123456",
	// })

}
