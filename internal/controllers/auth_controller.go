package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/dto"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/entity"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/responses"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/services"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := services.Auth{}
	userDto := dto.UserDto{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}
	var user entity.User

	validateLogin, msg := middleware.ValidateLogin(userDto)
	if !validateLogin {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: msg}})
	}

	err := services.UserCollection.FindOne(ctx, bson.M{"username": userDto.Username}).Decode(&user)
	if helper.ErrorIsNil(err) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgUsernameNotFound}})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDto.Password)); helper.ErrorIsNil(err) {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgLoginFail}})
	}

	token, err := auth.GenerateToken(user)
	if helper.ErrorIsNil(err) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.AccessToken: token}})
}

func RegisterAccount(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := services.Auth{}

	userDto := dto.UserDto{
		Username:        c.FormValue("username"),
		Email:           c.FormValue("email"),
		FirstName:       c.FormValue("first_name"),
		LastName:        c.FormValue("last_name"),
		Password:        c.FormValue("password"),
		Phone:           c.FormValue("phone"),
		ConfirmPassword: c.FormValue("confirm_password"),
	}

	validateRegister, msg := middleware.ValidateRegister(userDto)
	if !validateRegister {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: msg}})
	}

	newUser := auth.PutParamToRegisterUser(userDto)
	err := services.UserCollection.FindOne(ctx, bson.M{"username": newUser.Username}).Decode(&newUser)
	if !helper.ErrorIsNil(err) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgUsernameIsExist}})
	}

	result, err := services.UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgRegisterFail}})
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusCreated, Message: helper.Success, Data: &fiber.Map{helper.Data: result}})
}

func ForgotPassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := services.Auth{}

	code := auth.EncodeToString(6)
	username := c.FormValue("username")
	var user entity.User

	if helper.IsEmpty(username) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgInvalidUsername}})
	}

	err := services.UserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgUsernameNotFound}})
	}

	newPass, _ := bcrypt.GenerateFromPassword([]byte(code), 12)
	_, updatePass := services.UserCollection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"password": string(newPass)}})

	if helper.ErrorIsNil(updatePass) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgResetPasswordFail}})
	}

	sendNewPassword := auth.SendMail(user.Email, helper.SubjectResetPass, auth.ResetPassBodyContentSendMail(user.LastName, username, code))
	if !sendNewPassword {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgErrSendMail}})
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.Data: helper.MsgSendEmailForgotPassSuccess}})
}

func ChangePassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := services.Auth{}
	userInToken := c.Locals("user").(*jwt.Token)
	claims := userInToken.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	email := claims["email"].(string)
	lastName := claims["lastname"].(string)

	userDto := dto.UserDto{
		Password:        c.FormValue("password"),
		ConfirmPassword: c.FormValue("confirm_password"),
	}

	validateChangePass, msg := middleware.ValidateChangePass(userDto)
	if !validateChangePass {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: msg}})
	}

	passwordNew, _ := bcrypt.GenerateFromPassword([]byte(userDto.Password), 12)
	_, updatePass := services.UserCollection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"password": string(passwordNew)}})

	if helper.ErrorIsNil(updatePass) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgChangePassFail}})
	}

	mailChangePass := auth.SendMail(email, helper.SubjectChangePass, auth.ChangePassBodyContentSendMail(lastName, username))

	if !mailChangePass {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgErrSendMail}})
	}
	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.Data: helper.MsgChangePassSuccess}})
}
