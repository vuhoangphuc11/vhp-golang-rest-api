package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"log"
	"os"
)

func Middleware(app *fiber.App) {

	// Use middlewares for each route
	app.Use(
		helmet.New(), // add Helmet middleware
	)

	app.Use(basicauth.New(basicauth.Config{
		Realm: "Forbidden",
		Authorizer: func(user string, role string) bool {
			if user == "phucvu" {
				return true
			}
			if role == "Admin" {
				return true
			}
			return false
		},
	}))

	//vhp: Set config for CSRF middleware
	//csrfConfig := csrf.Config{
	//	KeyLookup:      "header:X-Csrf-Token", // string in the form of '<source>:<key>' that is used to extract token from the request
	//	CookieName:     "my_csrf_",            // name of the session cookie
	//	CookieSameSite: "Strict",              // indicates if CSRF cookie is requested by SameSite
	//	Expiration:     3 * time.Hour,         // expiration is the duration before CSRF token will expire
	//	KeyGenerator:   utils.UUID,            // creates a new CSRF token
	//}

	// Use middlewares for each route
	app.Use(
		csrf.New(), // add CSRF middleware with config
	)

	//vhp: Set config for Limiter middleware
	//limiterConfig := limiter.Config{
	//	Next: func(c *fiber.Ctx) bool {
	//		return c.IP() == "127.0.0.1" // limit will apply to this IP
	//	},
	//	Max:        20,                // max count of connections
	//	Expiration: 30 * time.Second,  // expiration time of the limit
	//	KeyGenerator: func(c *fiber.Ctx) string {
	//		return c.IP() // allows you to generate custom keys
	//	},
	//	LimitReached: func(c *fiber.Ctx) error {
	//		return c.SendStatus(fiber.StatusTooManyRequests) // called when a request hits the limit
	//	},
	//	SkipFailedRequests:     false,
	//	SkipSuccessfulRequests: false,
	//}

	// Use middlewares for each route
	app.Use(
		limiter.New(), // add Limiter middleware with config
	)

	//vhp: Define file to logs
	file, err := os.OpenFile("./vhp_logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	// Set config for logger
	loggerConfig := logger.Config{
		Output: file, // add file to save output
	}

	// Use middlewares for each route
	app.Use(
		logger.New(loggerConfig), // add Logger middleware with config
	)

}
