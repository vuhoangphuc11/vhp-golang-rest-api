package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/configs"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/routes"
)

func main() {
	app := fiber.New()

	//vhp: connect database
	configs.ConnectDB()

	//vhp: routes
	routes.UserRoute(app)

	//vhp: port
	err := app.Listen(":8087")
	if err != nil {
		return
	}

}
