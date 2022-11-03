package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/models"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/responses"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/services"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{helper.Data: err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := services.Validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{helper.Data: validationErr.Error()}})
	}

	if !helper.CheckPatternEmail(user.Email) {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgInValidFormatEmail}})
	}

	findUser := services.UserCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&user)

	if !helper.ErrorIsNil(findUser) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgUsernameIsExist}})
	}

	newUser := services.PutParamToCreateUser(user)

	result, err := services.UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.ResponseData{Status: http.StatusCreated, Message: helper.Success, Data: &fiber.Map{"data": result}})
}

func GetUserById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.Params("username")
	var user models.User

	err := services.UserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{"data": user}})
}

func UpdateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.Params("username")
	var user models.User

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{helper.Data: err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := services.Validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{helper.Data: validationErr.Error()}})
	}

	if !helper.CheckPatternEmail(user.Email) {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgInValidFormatEmail}})
	}

	updateUser := services.PutParamToUpdateUser(user)

	er := services.UserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if er != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgUsernameNotFound}})
	}

	_, err := services.UserCollection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": updateUser})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.Data: updateUser}})
}

func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	username := c.Params("username")
	defer cancel()

	result, err := services.UserCollection.DeleteOne(ctx, bson.M{"username": username})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgDeleteUserFail}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.ResponseData{Status: http.StatusNotFound, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgUsernameNotFound}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.Data: helper.MsgDeleteUserSuccess}},
	)
}

func GetAllUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User
	results, err := services.UserCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: err.Error()}})
		}
		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.Data: users}})
}

func ExportUserActive(c *fiber.Ctx) error {
	exportExcel := services.ExportExcel()
	if !exportExcel {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Error: helper.MsgExportExcelFail}})
	}
	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.Success: helper.MsgExportExcelSuccess}})

}

func HelloUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	return c.Status(http.StatusOK).JSON(responses.ResponseData{
		Status:  http.StatusOK,
		Message: helper.Success,
		Data:    &fiber.Map{"message": "Hello " + username},
	})
}

func Test1(c *fiber.Ctx) error {
	return c.SendString("Test ngon lanh`")
}
