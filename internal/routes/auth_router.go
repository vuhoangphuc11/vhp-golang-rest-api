package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/controllers"
)

func AuthRouter(app *fiber.App) {
	//vhp: Authenticate router
	app.Post("/api/auth/login", controllers.Login)
	app.Post("/api/auth/register", controllers.RegisterAccount)
	app.Post("/api/auth/forgot-password", controllers.ForgotPassword)

	//vhp: Unauthenticated route
	//app.Get("/", controllers.Accessible)

	//vhp: JWT Middleware
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	//vhp: Restricted router

}
