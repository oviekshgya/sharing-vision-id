package middlewares

import "github.com/gofiber/fiber/v2"

func NotFoundMiddleware(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendFile("./public/views/404.html")
}
