package database

import (
	"fmt"

	"github.com/RobertoSuarez/apialumni/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDB() {
	//dsn := "host=ec2-3-232-22-121.compute-1.amazonaws.com user=pocvdhmygcmcmp password=cbb38b84ecb8a348a6e99eca0864599e0c4410de56ca55a4bf94fb25e6226597 dbname=dbov41pluh9o22 port=5432 sslmode=disable"
	dsn := "postgres://pocvdhmygcmcmp:cbb38b84ecb8a348a6e99eca0864599e0c4410de56ca55a4bf94fb25e6226597@ec2-3-232-22-121.compute-1.amazonaws.com:5432/dbov41pluh9o22"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Database = db

	// Migraciones
	db.AutoMigrate(&models.Usuario{})
	db.AutoMigrate(&models.TipoUsuario{})
	db.AutoMigrate(&models.Admin{})
	db.AutoMigrate(&models.Alumni{})
	db.AutoMigrate(&models.OfertaLaboral{})

	// tipos de usuarios
	tipos := []models.TipoUsuario{}
	result := db.Find(&tipos)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	if len(tipos) == 0 {
		Database.Create(&models.TipoUsuario{Tipo: "admin"})
		Database.Create(&models.TipoUsuario{Tipo: "alumni"})
	}

	// db.Create(&models.Usuario{
	// 	Email:    "electrosonix12@gmail.com",
	// 	Password: "123456",
	// })

}
