package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/dto"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/entity"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/responses"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/services"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
	"net/http"
)

func AuthReq() func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Token is expired" {
				return c.Status(http.StatusUnauthorized).JSON(responses.ResponseData{Status: http.StatusUnauthorized, Message: helper.Error, Data: &fiber.Map{"message": "Token is expired!"}})
			}
			return c.Status(http.StatusUnauthorized).JSON(responses.ResponseData{Status: http.StatusUnauthorized, Message: helper.Error, Data: &fiber.Map{"message": "Unauthorized!"}})
		},
		SigningKey: []byte(services.SecretKey),
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

func ValidateCreateUser(c *fiber.Ctx, user entity.User) (bool, string) {
	//validate the request body
	if parserErr := c.BodyParser(&user); parserErr != nil {
		return false, parserErr.Error()
	}

	//use the validator library to validate required fields
	if validationErr := services.Validate.Struct(&user); validationErr != nil {
		return false, validationErr.Error()
	}

	if !helper.CheckPatternEmail(user.Email) {
		return false, helper.MsgInValidFormatEmail
	}
	return true, "success"
}

func ValidateChangePass(dto dto.UserDto) (bool, string) {
	if helper.IsEmpty(dto.Password) {
		return false, helper.MsgInValidPassword
	}
	if helper.IsEmpty(dto.ConfirmPassword) {
		return false, helper.MsgInValidConfirmPassword
	}
	if helper.NotMatch(dto.Password, dto.ConfirmPassword) {
		return false, helper.MsgInValidConfirmPasswordNotMatch
	}
	return true, helper.Success
}

func ValidateLogin(dto dto.UserDto) (bool, string) {
	if helper.IsEmpty(dto.Username) {
		return false, helper.MsgInvalidUsername
	}
	if helper.IsEmpty(dto.Password) {
		return false, helper.MsgInValidPassword
	}

	return true, helper.Success
}

func ValidateRegister(dto dto.UserDto) (bool, string) {
	if helper.IsEmpty(dto.Username) {
		return false, helper.MsgInvalidUsername
	}
	if helper.IsEmpty(dto.Email) {
		return false, helper.MsgInValidEmail
	}
	if !helper.CheckPatternEmail(dto.Email) {
		return false, helper.MsgInValidFormatEmail
	}
	if helper.IsEmpty(dto.FirstName) {
		return false, helper.MsgInValidFirstName
	}
	if helper.IsEmpty(dto.LastName) {
		return false, helper.MsgInValidLastName
	}
	if helper.IsEmpty(dto.Password) {
		return false, helper.MsgInValidPassword
	}
	if helper.IsEmpty(dto.Phone) {
		return false, helper.MsgInValidPhone
	}
	if helper.IsEmpty(dto.ConfirmPassword) {
		return false, helper.MsgInValidConfirmPassword
	}
	if helper.NotMatch(dto.Password, dto.ConfirmPassword) {
		return false, helper.MsgInValidConfirmPasswordNotMatch
	}
	return true, helper.Success
}
