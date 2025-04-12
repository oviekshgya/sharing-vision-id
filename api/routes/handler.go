package routes

import (
	"sharing-vision-id/api/controller"
	"sharing-vision-id/db"
)

var (
	UserController *controller.UserController
)

func InitialRoute() {
	UserController = controller.HandlerController(db.ConnDB)
}
