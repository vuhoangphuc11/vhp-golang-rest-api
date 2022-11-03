package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/models"
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

	paramArr := [2]string{}
	paramArr[0] = c.FormValue("username")
	paramArr[1] = c.FormValue("password")

	var user models.User

	validateLogin, msg := middleware.ValidateLogin(paramArr)
	if !validateLogin {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: msg}})
	}

	err := services.UserCollection.FindOne(ctx, bson.M{"username": paramArr[0]}).Decode(&user)
	if helper.ErrorIsNil(err) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgUsernameNotFound}})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(paramArr[1])); helper.ErrorIsNil(err) {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgLoginFail}})
	}

	token, err := services.GenerateToken(user)
	if helper.ErrorIsNil(err) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.AccessToken: token}})
}

func RegisterAccount(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramArr := [7]string{}
	paramArr[0] = c.FormValue("username")
	paramArr[1] = c.FormValue("email")
	paramArr[2] = c.FormValue("first_name")
	paramArr[3] = c.FormValue("last_name")
	paramArr[4] = c.FormValue("password")
	paramArr[5] = c.FormValue("phone")
	paramArr[6] = c.FormValue("confirm_password")

	validateRegister, msg := middleware.ValidateRegister(paramArr)
	if !validateRegister {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: msg}})
	}

	newUser := services.PutParamToRegisterUser(paramArr)
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

	code := services.EncodeToString(6)
	username := c.FormValue("username")
	var user models.User

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

	sendNewPassword := services.SendMail(user.Email, helper.SubjectResetPass, services.ResetPassBodyContentSendMail(user.LastName, username, code))
	if !sendNewPassword {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgErrSendMail}})
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.Data: helper.MsgSendEmailForgotPassSuccess}})
}

func ChangePassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userInToken := c.Locals("user").(*jwt.Token)
	claims := userInToken.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	email := claims["email"].(string)
	lastName := claims["lastname"].(string)

	paramArr := [2]string{}
	paramArr[0] = c.FormValue("password")
	paramArr[1] = c.FormValue("confirm_password")

	validateChangePass, msg := middleware.ValidateChangePass(paramArr)
	if !validateChangePass {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: msg}})
	}

	passwordNew, _ := bcrypt.GenerateFromPassword([]byte(paramArr[0]), 12)
	_, updatePass := services.UserCollection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"password": string(passwordNew)}})

	if helper.ErrorIsNil(updatePass) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgChangePassFail}})
	}

	mailChangePass := services.SendMail(email, helper.SubjectChangePass, services.ChangePassBodyContentSendMail(lastName, username))
	if !mailChangePass {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{helper.Data: helper.MsgErrSendMail}})
	}
	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.Data: helper.MsgChangePassSuccess}})
}
