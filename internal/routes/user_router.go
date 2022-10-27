package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/controllers"
)

func UserRoute(app *fiber.App) {
	//vhp: User router
	app.Get("/api/user/hello", controllers.AuthReq(), controllers.HelloUser)
	app.Get("/api/user/get-all", controllers.GetAllUser)
	app.Get("/api/user/get-user/:userId", controllers.GetUserById)

	//vhp: JWT Middleware
	//app.Use(jwtware.New(jwtware.Config{
	//	SigningKey: []byte("secret"),
	//}))

	app.Post("/api/user/create-user", controllers.AuthReq(), controllers.CreateUser)
	app.Put("/api/user/update-user/:userId", controllers.AuthReq(), controllers.UpdateUser)
	app.Delete("/api/user/delete-user/:userId", controllers.AuthReq(), controllers.DeleteUser)
}
