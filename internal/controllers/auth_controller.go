package controllers

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/models"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/responses"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/mail"
	"net/smtp"
	"time"
)

const SecretKey = "WuH0aNqFuc"

func Login(c *fiber.Ctx) error {
	currentTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	username := c.FormValue("username")
	password := c.FormValue("password")
	var user models.User
	defer cancel()

	if helper.IsEmpty(username) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": "Please enter your username!"}})
	}
	if helper.IsEmpty(password) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": "Please enter your password!"}})
	}

	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if helper.ErrorIsNil(err) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{"data": "Username not found!"}})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); helper.ErrorIsNil(err) {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{"data": "Login failed, please try again!"}})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    user.Username,
		ExpiresAt: &jwt.NumericDate{Time: currentTime.Add(time.Hour * 6)},
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if helper.ErrorIsNil(err) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{helper.AccessToken: token}})
}

func RegisterAccount(c *fiber.Ctx) error {
	var createAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	username := c.FormValue("username")
	email := c.FormValue("email")
	firstName := c.FormValue("first_name")
	lastName := c.FormValue("last_name")
	password := c.FormValue("password")
	phone := c.FormValue("phone")
	confirmPassword := c.FormValue("confirm_password")

	var user models.User
	defer cancel()

	if helper.IsEmpty(username) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Please enter your username!"}})
	}
	if helper.IsEmpty(email) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Please enter your email!"}})
	}
	if !checkPatternEmail(email) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Invalid email, please try again!"}})
	}
	if helper.IsEmpty(password) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Please enter your password!"}})
	}
	if helper.IsEmpty(confirmPassword) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Please enter your confirm password!"}})
	}
	if helper.IsEmpty(firstName) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Please enter your first name!"}})
	}
	if helper.IsEmpty(lastName) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Please enter your last name!"}})
	}
	if helper.IsEmpty(phone) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Please enter your phone!!"}})
	}
	if helper.NotMatch(password, confirmPassword) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Password and confirm Password not match!"}})
	}

	passwordNew, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	newUser := models.User{
		Id:        primitive.NewObjectID(),
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(passwordNew),
		Phone:     phone,
		IsActive:  true,
		Role:      "User",
		CreateAt:  createAt,
	}

	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if !helper.ErrorIsNil(err) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Account already exists!"}})
	}

	result, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Register fail!"}})
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusCreated, Message: helper.Success, Data: &fiber.Map{helper.Data: result}})

}

func ForgotPassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	code := encodeToString(6)
	username := c.FormValue("username")
	var user models.User
	defer cancel()

	if helper.IsEmpty(username) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Please enter your username!"}})
	}

	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Account not exists!"}})
	}

	newPass, _ := bcrypt.GenerateFromPassword([]byte(code), 12)

	_, updatePass := userCollection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"password": string(newPass)}})

	if helper.ErrorIsNil(updatePass) {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResponseData{Status: http.StatusInternalServerError, Message: helper.Error, Data: &fiber.Map{helper.Data: "Reset password fail!"}})
	}

	// Sender data.
	from := "phucvhps12860@fpt.edu.vn"
	password := "pbqrhwveiavkjvyq"
	subject := "[Reset password by VHP]"
	body := "Reset password for account " + username + "." +
		"\nNew password: " + code +
		"\n\n\n" +
		"Vu Hoang Phuc\n" +
		"Golang Developer\n" +
		"Best regards."

	// Receiver email address.
	to := []string{
		user.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += fmt.Sprintf("\r\n%s\r\n", body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	sendMail := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if sendMail != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ResponseData{Status: http.StatusBadRequest, Message: helper.Error, Data: &fiber.Map{"data": sendMail.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.ResponseData{Status: http.StatusOK, Message: helper.Success, Data: &fiber.Map{"data": "Sent an email reset the password"}})
}

func encodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		_ = fmt.Errorf(err.Error())
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func checkPatternEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func Accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)
	return c.SendString("Welcome " + name)
}
