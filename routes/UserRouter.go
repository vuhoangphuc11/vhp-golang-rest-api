package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/controllers"
)

func UserRoute(app *fiber.App) {
	app.Get("/api/user", controllers.GetAllUser)
	app.Get("/api/user/get-all", controllers.GetAllUser)
	app.Get("/api/user/get-user/:userId", controllers.GetUserById)
	app.Post("/api/user/create-user", controllers.CreateUser)
	app.Put("/api/user/update-user/:userId", controllers.UpdateUser)
	app.Delete("/api/user/delete-user/:userId", controllers.DeleteUser)
}
