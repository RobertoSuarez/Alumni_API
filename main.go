package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("port", "3000")

	app := fiber.New()

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

	log.Println("listen to :" + viper.GetString("port"))
	app.Listen(":" + viper.GetString("port"))
}
