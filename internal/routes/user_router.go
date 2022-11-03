package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/controllers"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/middleware"
)

func UserRoute(app *fiber.App) {
	//vhp: api version 1
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/user")

	//vhp: user api
	user.Get("/hello", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.HelloUser)
	user.Get("/get-all", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.GetAllUser)
	user.Get("/get-user/:username", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.GetUserById)
	user.Post("/create-user", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager), controllers.CreateUser)
	user.Put("/update-user/:username", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.UpdateUser)
	user.Delete("/delete-user/:username", middleware.AuthReq(), middleware.AuthorReq(helper.Admin), controllers.DeleteUser)
	user.Get("/export-user-active", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager), controllers.ExportUserActive)
	user.Post("/test", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.Test1)
}
