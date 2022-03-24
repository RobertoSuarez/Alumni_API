package main

import (
	"fmt"
	"log"

	"github.com/RobertoSuarez/apialumni/config"
	"github.com/RobertoSuarez/apialumni/controllers"
	"github.com/RobertoSuarez/apialumni/database"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logrus.Trace("API REST Alumni")
}

func main() {
	configvar := viper.New()
	configvar.AddConfigPath(".")
	configvar.SetConfigName("app")
	configvar.SetConfigType("env")

	configvar.AutomaticEnv()

	if err := configvar.ReadInConfig(); err != nil {
		fmt.Println("Error al leer las variables de configuración")
		log.Println(err)
	} else {
		fmt.Println("Las variables se establecierón correctamente")
	}

	database.ConnectDB(configvar)

	app := fiber.New()

	app.Use(cors.New())
	api := app.Group("/api/v1")

	config.Use(api.Group("/auth"), controllers.NewControllerAuth())
	config.Use(api.Group("/users"), controllers.NewControllerUsuario())
	config.Use(api.Group("/ofertas"), controllers.NewControllerOfertaLaboral())
	config.Use(api.Group("/educacion"), controllers.NewControllerEducacion())

	logrus.Info("listen to :" + viper.GetString("port"))
	app.Listen(":" + viper.GetString("port"))
}
