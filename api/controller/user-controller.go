package controller

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sharing-vision-id/internal/models"
	"sharing-vision-id/internal/service"
	"sharing-vision-id/pkg"
	"strconv"
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

func (controller *UserController) Create(c *fiber.Ctx) error {
	responseInitial := pkg.InitialResponse{Ctx: c}
	var input models.Post
	if err := c.BodyParser(&input); err != nil {
		return responseInitial.Respose(fiber.StatusBadRequest, "invalid payload request: "+err.Error(), true, nil)
	}

	result, err := controller.UserService.CreatePost(input)
	if err != nil {
		return responseInitial.Respose(fiber.StatusInternalServerError, "error creating post: "+err.Error(), true, nil)
	}

	return responseInitial.Respose(http.StatusCreated, "Created", false, result)
}

func (controller *UserController) GetAll(c *fiber.Ctx) error {
	responseInitial := pkg.InitialResponse{Ctx: c}
	pageSize, _ := strconv.Atoi(c.Params(":limit"))
	page, _ := strconv.Atoi(c.Params(":offset"))
	result, err := controller.UserService.GetData(0, page, pageSize)
	if err != nil {
		return responseInitial.Respose(fiber.StatusInternalServerError, "error creating post: "+err.Error(), true, nil)
	}

	return responseInitial.Respose(http.StatusOK, "Created", false, result)
}

func (controller *UserController) GetById(c *fiber.Ctx) error {
	responseInitial := pkg.InitialResponse{Ctx: c}
	id, _ := strconv.Atoi(c.Params(":id"))
	result, err := controller.UserService.GetData(id, 0, 0)
	if err != nil {
		return responseInitial.Respose(fiber.StatusInternalServerError, "error creating post: "+err.Error(), true, nil)
	}

	return responseInitial.Respose(http.StatusOK, "Created", false, result)
}

func (controller *UserController) Update(c *fiber.Ctx) error {
	responseInitial := pkg.InitialResponse{Ctx: c}
	var input models.Post
	if err := c.BodyParser(&input); err != nil {
		return responseInitial.Respose(fiber.StatusBadRequest, "invalid payload request: "+err.Error(), true, nil)
	}

	result, err := controller.UserService.Update(input)
	if err != nil {
		return responseInitial.Respose(fiber.StatusInternalServerError, "error creating post: "+err.Error(), true, nil)
	}

	return responseInitial.Respose(http.StatusAccepted, "update", false, result)
}

func (controller *UserController) Delete(c *fiber.Ctx) error {
	responseInitial := pkg.InitialResponse{Ctx: c}

	id, _ := strconv.Atoi(c.Params(":id"))
	result, err := controller.UserService.Delete(uint(id))
	if err != nil {
		return responseInitial.Respose(fiber.StatusInternalServerError, "error creating post: "+err.Error(), true, nil)
	}

	return responseInitial.Respose(http.StatusAccepted, "update", false, result)
}
