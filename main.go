package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("port", "3000")

	app := fiber.New()

	app.Static("/*", "./dist")

	log.Println("listen to :" + viper.GetString("port"))
	app.Listen(":" + viper.GetString("port"))
}
