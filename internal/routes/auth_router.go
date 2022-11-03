package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/controllers"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/middleware"
)

func AuthRouter(app *fiber.App) {
	//vhp: api version 1
	api := app.Group("/api")
	v1 := api.Group("/v1")
	auth := v1.Group("/auth")

	//vhp: authenticate api
	auth.Post("/login", controllers.Login)
	auth.Post("/register", controllers.RegisterAccount)
	auth.Post("/forgot-password", controllers.ForgotPassword)
	auth.Post("/change-password", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.ChangePassword)
}
