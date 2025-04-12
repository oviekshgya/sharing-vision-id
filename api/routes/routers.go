package routes

import (
	"github.com/gofiber/fiber/v2"
	middlewares "sharing-vision-id/api/middleware"
)

var (
	Router *fiber.App
)

func Route() {
	article := Router.Group("/article")
	article.Use(middlewares.CORSMiddleware())
	{
		article.Post("/", UserController.Create)
		article.Get("/:limit/:offet", UserController.GetAll)
		article.Get("/:id", UserController.GetById)
		article.Put("/:id", UserController.Update)
		article.Delete("/:id", UserController.Delete)
	}

}
