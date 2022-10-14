package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/configs"
	routes2 "github.com/vuhoangphuc11/vhp-golang-rest-api/routes"
)

func main() {
	app := fiber.New()

	//vhp: connect database
	configs.ConnectDB()

	//vhp: routes
	routes2.UserRoute(app)
	routes2.AuthRouter(app)

	//vhp: port
	err := app.Listen(":8087")
	if err != nil {
		return
	}
}
