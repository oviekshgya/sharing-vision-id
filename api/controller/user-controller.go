package controller

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sharing-vision-id/internal/service"
	"sharing-vision-id/pkg"
)

type UserController struct {
	DB          *gorm.DB
	Client      *mongo.Client
	UserService service.UserService
}

func HandlerController(db *gorm.DB) *UserController {
	if db == nil {
		log.Println("Database [HandlerController] connection is nil")
	}

	return &UserController{
		DB: db,
		UserService: service.UserService{
			DB: db,
		},
	}
}

func (controller *UserController) Welcome(c *fiber.Ctx) error {
	responseInitial := pkg.InitialResponse{Ctx: c}
	return responseInitial.Respose(http.StatusOK, "Welcome", false, map[string]interface{}{
		"message": "Welcome Shagya-Tech Payment" + pkg.GoogleOAuthConfig.ClientID,
	})
}
