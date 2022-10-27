package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/configs"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/routes"
	"log"
	"os"
	"time"
)

func main() {
	app := fiber.New()

	//vhp: Define file to logs
	file, er := os.OpenFile("./vhp_logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if er != nil {
		log.Fatalf("error opening file: %v", er)
	}
	defer file.Close()

	//vhp: Set config for logger
	loggerConfig := logger.Config{
		Output: file, // add file to save output
	}

	//vhp: Use middlewares for each route
	app.Use(
		logger.New(loggerConfig), // add Logger middleware with config
	)

	//vhp: limit req
	limit := limiter.Config{
		Max:        1,
		Expiration: time.Second,
	}

	app.Use(limiter.New(limit))

	//vhp: connect database
	configs.ConnectDB()

	//vhp: routes
	routes.UserRoute(app)
	routes.AuthRouter(app)

	//vhp: port
	err := app.Listen(":8087")
	if err != nil {
		return
	}
}
