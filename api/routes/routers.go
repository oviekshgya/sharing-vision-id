package routes

import (
	"github.com/gofiber/fiber/v2"
)

var (
	Router *fiber.App
)

func Route() {
	Router.Static("/demo", "./public/views")

}
