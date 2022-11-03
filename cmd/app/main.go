package main

import (
	"github.com/gkampitakis/fiber-modules/gracefulshutdown"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	file, er := os.OpenFile(os.Getenv("LOG_FILE_NAME"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	app.Use(cors.New())

	//vhp: limit req
	limit := limiter.Config{
		Max:        1,
		Expiration: time.Second,
	}

	app.Use(limiter.New(limit))

	//vhp: connect database
	configs.ConnectDB()

	//vhp: use routes
	routes.UserRoute(app)
	routes.AuthRouter(app)

	//vhp: port
	//err := app.Listen(os.Getenv("PORT"))
	//if err != nil {
	//	log.Panic(err)
	//}

	//vhp: Registers graceful shutdown with default config
	gracefulshutdown.Listen(app, os.Getenv("PORT"), gracefulshutdown.Default())

}
