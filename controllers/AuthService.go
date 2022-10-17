package controllers

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/models"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/responses"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/smtp"
	"time"
)

const SecretKey = "WuH0aNqFuc"

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	username := c.FormValue("username")
	password := c.FormValue("password")
	var user models.User
	defer cancel()

	if util.IsEmpty(username) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{"data": "Please enter your username!"}})
	}
	if util.IsEmpty(password) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{"data": "Please enter your password!"}})
	}

	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if util.ErrorIsNil(err) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{"data": "Username not found!"}})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); util.ErrorIsNil(err) {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: util.Error, Data: &fiber.Map{"data": "Login failed, please try again!"}})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Username,
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if util.ErrorIsNil(err) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: util.Success, Data: &fiber.Map{util.AccessToken: token}})
}

func RegisterAccount(c *fiber.Ctx) error {
	var createAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	username := c.FormValue("username")
	email := c.FormValue("email")
	firstName := c.FormValue("firstname")
	lastName := c.FormValue("lastname")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	var user models.User
	defer cancel()

	if util.IsEmpty(username) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Please enter your username!"}})
	}
	if util.IsEmpty(email) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Please enter your email!"}})
	}
	if util.IsEmpty(password) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Please enter your password!"}})
	}
	if util.IsEmpty(confirmPassword) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Please enter your confirm password!"}})
	}
	if util.IsEmpty(firstName) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Please enter your first name!"}})
	}
	if util.IsEmpty(lastName) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Please enter your last name!"}})
	}
	if util.NotMatch(password, confirmPassword) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Password and confirm Password not match!"}})
	}

	passwordNew, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	newUser := models.User{
		Id:        primitive.NewObjectID(),
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(passwordNew),
		IsActive:  true,
		Role:      "User",
		CreateAt:  createAt,
	}

	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err == nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Account already exists!"}})
	}

	result, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Register fail!"}})
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusCreated, Message: util.Success, Data: &fiber.Map{util.Data: result}})

}

func ForgotPassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	code := EncodeToString(6)
	username := c.FormValue("username")
	email := c.FormValue("email")
	var user models.User
	defer cancel()

	if util.IsEmpty(username) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Please enter your username!"}})
	}
	if util.IsEmpty(email) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Please enter your email!"}})
	}

	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: util.Error, Data: &fiber.Map{util.Data: "Account not exists!"}})
	}

	// Sender data.
	from := "phucvhps12860@fpt.edu.vn"
	password := "pbqrhwveiavkjvyq"
	subject := "[Reset password for user by VHP]"
	body := "Reset password for account " + username + ".\nNew Pass: " + code

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	//message := []byte("Reset password for account " + username + " is " + code)

	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += fmt.Sprintf("\r\n%s\r\n", body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	sendError := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if sendError != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: util.Error, Data: &fiber.Map{"data": sendError.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: util.Success, Data: &fiber.Map{"data": "Sent an email reset the password"}})
}

func Accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		fmt.Errorf(err.Error())
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
