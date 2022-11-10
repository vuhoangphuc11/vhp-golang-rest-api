package services

import (
	"crypto/rand"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/configs"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/dto"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/entity"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/smtp"
	"os"
	"time"
)

var UserCollection *mongo.Collection = configs.GetCollection(configs.DB, "user")
var Validate = validator.New()

const SecretKey = "WuH0aNqFuc"

type AuthInterface interface {
	ResetPassBodyContentSendMail(name, username, code string) string
	ChangePassBodyContentSendMail(lastName, username string) string
	PutParamToRegisterUser(dto dto.UserDto) entity.User
	GenerateToken(user entity.User) (string, error)
	SendMail(email, subject, body string) bool
	EncodeToString(max int) string
}

type Auth struct {
}

func (r *Auth) ResetPassBodyContentSendMail(name, username, code string) string {
	body :=
		"\nHi " + name + "," +
			"\n\nYour account's password is " + username + "." +
			"\nHas been changed to: " + code +
			"\n\n\n" +
			helper.SignatureMail

	return body
}

func (r *Auth) ChangePassBodyContentSendMail(lastName, username string) string {
	body :=
		"\nHi " + lastName + "," +
			"\n\nPassword of account: " + username + "." +
			"\nJust changed" +
			"\n\n\n" +
			helper.SignatureMail
	return body
}

func (r *Auth) GenerateToken(user entity.User) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"username": user.Username,
		"lastname": user.LastName,
		"role":     user.Role,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	}
	// Create token
	createToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	token, err := createToken.SignedString([]byte(SecretKey))
	return token, err
}

func (r *Auth) PutParamToRegisterUser(dto dto.UserDto) entity.User {
	var createAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	passwordNew, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), 12)

	newUser := entity.User{
		Id:        primitive.NewObjectID(),
		Username:  dto.Username,
		Email:     dto.Email,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Password:  string(passwordNew),
		Phone:     dto.Phone,
		IsActive:  true,
		Role:      helper.User,
		CreateAt:  createAt,
	}
	return newUser
}

func (r *Auth) SendMail(email, subject, body string) bool {
	// Sender data.
	from := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Prepare Content
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += fmt.Sprintf("\r\n%s\r\n", body)

	// Authentication.
	authentication := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	sendMail := smtp.SendMail(smtpHost+":"+smtpPort, authentication, from, to, []byte(message))

	if sendMail != nil {
		return false
	}
	return true
}

func (r *Auth) EncodeToString(max int) string {
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
