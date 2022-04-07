package main

import (
	"fmt"
	"log"

	"github.com/RobertoSuarez/apialumni/config"
	"github.com/RobertoSuarez/apialumni/controllers"
	"github.com/RobertoSuarez/apialumni/models"
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

	// Iniciamos la base de datos con las tablas
	models.InitDataBaseTable(configvar)

	app := fiber.New()

	app.Use(cors.New())
	api := app.Group("/api/v1")

	config.UseMount("/auth", api, controllers.NewControllerAuth())
	config.UseMount("/usuarios", api, controllers.NewControllerUsuario())
	config.UseMount("/empleos", api, controllers.NewEmpleo())
	config.UseMount("/educacion", api, controllers.NewControllerEducacion())
	config.UseMount("/grupos", api, controllers.NewGrupo())
	config.UseMount("/empresas", api, controllers.NewEmpresa())
	config.UseMount("/areas", api, controllers.NewControllerArea())

	logrus.Info("listen to :" + configvar.GetString("port"))
	app.Listen(":" + configvar.GetString("port"))
}
