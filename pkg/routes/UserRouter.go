package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/services"
)

func UserRoute(app *fiber.App) {
	app.Get("/api/user", services.GetAllUser)
	app.Get("/api/user/get-all", services.GetAllUser)
	app.Get("/api/user/get-user/:userId", services.GetUserById)
	app.Post("/api/user/create-user", services.CreateUser)
	app.Put("/api/user/update-user/:userId", services.UpdateUser)
	app.Delete("/api/user/delete-user/:userId", services.DeleteUser)
}
