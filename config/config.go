package config

import "github.com/gofiber/fiber/v2"

type ConfigMicroServicio interface {
	ConfigPath(router fiber.Router)
}

func Use(r fiber.Router, microServicio ConfigMicroServicio) {
	microServicio.ConfigPath(r)
}
