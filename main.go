package main

import (
	"time"

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
	viper.AutomaticEnv()
	viper.SetDefault("port", "3000")

	database.ConnectDB()

	app := fiber.New()

	app.Use(cors.New())
	api := app.Group("/api/v1")

	config.Use(api.Group("/auth"), controllers.NewControllerAuth())
	config.Use(api.Group("/users"), controllers.NewControllerUsuario())
	config.Use(api.Group("/ofertas"), controllers.NewControllerOfertaLaboral())
	config.Use(api.Group("/educacion"), controllers.NewControllerEducacion())

	// Frontend
	app.Static("/", "./dist", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})
	app.Static("/*", "./dist", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	logrus.Info("listen to :" + viper.GetString("port"))
	app.Listen(":" + viper.GetString("port"))
}
