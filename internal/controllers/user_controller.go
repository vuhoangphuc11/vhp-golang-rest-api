package controllers

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/configs"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/models"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/responses"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "user")
var validate = validator.New()

const (
	SheetName = "User List"
)

func CreateUser(c *fiber.Ctx) error {
	var createAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{"data": validationErr.Error()}})
	}

	if !checkPatternEmail(user.Email) {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": "Invalid email, please try again!"}})
	}

	findUser := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&user)

	if !helper.ErrorIsNil(findUser) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Username already exists!"}})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	newUser := models.User{
		Id:        primitive.NewObjectID(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  string(password),
		Age:       user.Age,
		Gender:    user.Gender,
		Phone:     user.Phone,
		IsActive:  user.IsActive,
		Role:      user.Role,
		CreateAt:  createAt,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.ResponseData{Status: http.StatusCreated, Message: helper.Success, Data: &fiber.Map{"data": result}})
}

func GetUserById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{"data": user}})
}

func UpdateUser(c *fiber.Ctx) error {
	var updateAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{"data": validationErr.Error()}})
	}
	encryptPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	updateUser := bson.M{
		"email":     user.Email,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"password":  string(encryptPass),
		"age":       user.Age,
		"gender":    user.Gender,
		"phone":     user.Phone,
		"isactive":  user.IsActive,
		"role":      user.Role,
		"updateat":  updateAt,
	}

	if !checkPatternEmail(user.Email) {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": "Invalid email, please try again!"}})
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": updateUser})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{"data": updatedUser}})
}

func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.ResponseData{Status: http.StatusNotFound, Message: helper.Error, Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{"data": users}},
	)
}

func ExportUserActive(c *fiber.Ctx) error {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", SheetName)

	titleStyle, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 28, Color: "2B4492", Bold: true}})
	err = f.MergeCell(SheetName, "B2", "E2")
	err = f.SetCellStyle(SheetName, "B2", "E2", titleStyle)
	err = f.SetSheetRow(SheetName, "B2", &[]interface{}{"User List Active"})

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 13, Bold: true, Color: "2B4492"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	err = f.SetCellStyle(SheetName, "B6", "K6", headerStyle)
	err = f.SetSheetRow(SheetName, "B6", &[]interface{}{"STT", "Username", "FullName", "Email", "Role", "Gender", "Phone", "Age", "Active", "Note"})

	listUserActive := GetListUserIsActive()

	var fillColor string
	//for j, _ := range listUserActive {
	//
	//}

	for i, v := range listUserActive {
		if i%2 == 0 {
			fillColor = "F3F3F3"
		} else {
			fillColor = "FFFFFF"
		}
		bodyStyle, _ := f.NewStyle(&excelize.Style{
			Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{fillColor}},
			Font:      &excelize.Font{Color: "666666"},
			Alignment: &excelize.Alignment{Vertical: "left"},
		})
		err = f.SetCellStyle(SheetName, fmt.Sprintf("B%d", i+7), fmt.Sprintf("K%d", i+7), bodyStyle)
		//err = f.SetSheetRow(SheetName, fmt.Sprintf(" B%d", i+5), &[]interface{}{v.Username, v.LastName + " " + v.FirstName, v.Email, v.Role, v.Gender, v.Phone, v.Age, v.IsActive, " "})
		f.SetCellValue(SheetName, fmt.Sprintf("B%d", i+7), i+1)
		f.SetCellValue(SheetName, fmt.Sprintf("C%d", i+7), v.Username)
		f.SetCellValue(SheetName, fmt.Sprintf("D%d", i+7), v.LastName+" "+v.FirstName)
		f.SetCellValue(SheetName, fmt.Sprintf("E%d", i+7), v.Email)
		f.SetCellValue(SheetName, fmt.Sprintf("F%d", i+7), v.Role)
		f.SetCellValue(SheetName, fmt.Sprintf("G%d", i+7), v.Gender)
		f.SetCellValue(SheetName, fmt.Sprintf("H%d", i+7), v.Phone)
		f.SetCellValue(SheetName, fmt.Sprintf("I%d", i+7), v.Age)
		f.SetCellValue(SheetName, fmt.Sprintf("J%d", i+7), v.IsActive)
		f.SetCellValue(SheetName, fmt.Sprintf("K%d", i+7), "")
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": "Style error"}})
	}

	if err := f.SaveAs("simple.xlsx"); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{"data": "Export successfully!"}},
	)
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

func GetListUserIsActive() []models.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.D{{"isactive", bson.M{"$exists": true}}})

	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return nil
		}

		users = append(users, singleUser)
	}

	return users
}

func Test1(c *fiber.Ctx) error {
	return c.SendString("Test ngon lanh`")
}
