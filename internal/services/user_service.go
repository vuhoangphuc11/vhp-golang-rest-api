package services

import (
	"context"
	"fmt"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/internal/models"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const SheetName = "User List"

func PutParamToCreateUser(user models.User) models.User {
	var createAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
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

	return newUser
}

func PutParamToUpdateUser(user models.User) bson.M {
	var updateAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
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

	return updateUser
}

func ExportExcel() bool {
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
		err = f.SetCellValue(SheetName, fmt.Sprintf("B%d", i+7), i+1)
		err = f.SetCellValue(SheetName, fmt.Sprintf("C%d", i+7), v.Username)
		err = f.SetCellValue(SheetName, fmt.Sprintf("D%d", i+7), v.LastName+" "+v.FirstName)
		err = f.SetCellValue(SheetName, fmt.Sprintf("E%d", i+7), v.Email)
		err = f.SetCellValue(SheetName, fmt.Sprintf("F%d", i+7), v.Role)
		err = f.SetCellValue(SheetName, fmt.Sprintf("G%d", i+7), v.Gender)
		err = f.SetCellValue(SheetName, fmt.Sprintf("H%d", i+7), v.Phone)
		err = f.SetCellValue(SheetName, fmt.Sprintf("I%d", i+7), v.Age)
		err = f.SetCellValue(SheetName, fmt.Sprintf("J%d", i+7), v.IsActive)
		err = f.SetCellValue(SheetName, fmt.Sprintf("K%d", i+7), "")
	}

	if err != nil {
		return false
	}

	if err := f.SaveAs("simple.xlsx"); err != nil {
		return false
	}

	return true
}

func GetListUserIsActive() []models.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := UserCollection.Find(ctx, bson.D{{"isactive", bson.M{"$exists": true}}})

	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return nil
		}

		users = append(users, singleUser)
	}

	return users
}
