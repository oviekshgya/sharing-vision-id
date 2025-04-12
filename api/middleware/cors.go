package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"sharing-vision-id/db"
	"sharing-vision-id/internal/models"
)

func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "https://shagyaapi-production.up.railway.app", //Deployment
		//AllowOrigins:     "http://127.0.0.1", //Deployment
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-KEY, X-SIGNATURE, X-TIMESTAMP",
		AllowCredentials: true,
	})
}

//func BasicAuthMiddleware() fiber.Handler {
//	return basicauth.New(basicauth.Config{
//		Users: map[string]string{
//			pkg.USERNAME: pkg.PASSWORD,
//		},
//		Unauthorized: func(c *fiber.Ctx) error {
//			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
//				"error": "Unauthorized",
//			})
//		},
//	})
//}

func BasicAuthMiddleware() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Authorizer: func(username, password string) bool {
			var client = models.ClientImpl{
				DB: db.DBMongo,
			}
			user, err := client.FindOne(username)
			if err != nil {
				return false
			}

			return user.PasswordAuth == password
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		},
	})
}

func APIKeyMiddleware(validAPIKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-KEY")

		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API Key is missing",
			})
		}

		if apiKey != validAPIKey {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid API Key",
			})
		}

		return c.Next()
	}
}
