package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/controllers"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/middleware"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
)

func UserRoute(app *fiber.App) {
	//vhp: User router
	app.Get("/api/user/hello", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.HelloUser)
	app.Get("/api/user/get-all", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.GetAllUser)
	app.Get("/api/user/get-user/:userId", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.GetUserById)
	app.Post("/api/user/create-user", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager), controllers.CreateUser)
	app.Put("/api/user/update-user/:userId", middleware.AuthReq(), middleware.AuthorReq(helper.Admin, helper.Manager, helper.User), controllers.UpdateUser)
	app.Delete("/api/user/delete-user/:userId", middleware.AuthReq(), middleware.AuthorReq(helper.Admin), controllers.DeleteUser)
	app.Post("/api/user/test", middleware.AuthReq(), middleware.AuthorReq(helper.Admin), controllers.Test1)
}
