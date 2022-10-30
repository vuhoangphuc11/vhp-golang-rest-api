package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/controllers"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/responses"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
	"net/http"
)

func AuthReq() func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusUnauthorized).JSON(responses.ResponseData{Status: http.StatusUnauthorized, Message: helper.Error, Data: &fiber.Map{"message": "Unauthorized!"}})
		},
		SigningKey: []byte(controllers.SecretKey),
	})
}

func AuthorReq(listRole ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)

		for _, v := range listRole {
			if v == role {
				return c.Next()
			}
		}
		return c.Status(http.StatusUnauthorized).JSON(responses.ResponseData{Status: http.StatusUnauthorized, Message: helper.Error, Data: &fiber.Map{"message": "You don't have permission!"}})
	}
}
