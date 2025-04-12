package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "https://shagya-tech.my.id", //Deployment
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-KEY, X-SIGNATURE, X-TIMESTAMP",
		AllowCredentials: true,
	})
}
