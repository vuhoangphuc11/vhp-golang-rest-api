package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/controllers"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/middleware"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
)

func AuthRouter(app *fiber.App) {
	//vhp: Authenticate router
	app.Post("/api/auth/login", controllers.Login)
	app.Post("/api/auth/register", controllers.RegisterAccount)
	app.Post("/api/auth/forgot-password", controllers.ForgotPassword)
	app.Post("/api/auth/change-password", middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.ChangePassword())
}
