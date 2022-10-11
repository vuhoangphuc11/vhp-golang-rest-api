package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/data"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/dto"
	"net/http"
	"strconv"
)

func GetAllUser(writer http.ResponseWriter, request *http.Request) {
	responseWithJson(writer, http.StatusOK, data.Users)
}
func GetUserById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		responseWithJson(writer, http.StatusBadRequest, map[string]string{"message": "Invalid uesr id"})
		return
	}

	for _, user := range data.Users {
		if user.ID == id {
			responseWithJson(writer, http.StatusOK, user)
			return
		}
	}

	responseWithJson(writer, http.StatusNotFound, map[string]string{"message": "User not found"})
}

func CreateUser(writer http.ResponseWriter, request *http.Request) {
	var newUser dto.User
	if err := json.NewDecoder(request.Body).Decode(&newUser); err != nil {
		responseWithJson(writer, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	newUser.ID = generateId(data.Users)
	data.Users = append(data.Users, newUser)

	responseWithJson(writer, http.StatusCreated, newUser)
}

func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		responseWithJson(writer, http.StatusBadRequest, map[string]string{"message": "Invalid user id"})
		return
	}

	var updateUser dto.User
	if err := json.NewDecoder(request.Body).Decode(&updateUser); err != nil {
		responseWithJson(writer, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}
	updateUser.ID = id

	for i, user := range data.Users {
		if user.ID == id {
			data.Users[i] = updateUser
			responseWithJson(writer, http.StatusOK, updateUser)
			return
		}
	}

	responseWithJson(writer, http.StatusNotFound, map[string]string{"message": "User not found"})
}
func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		responseWithJson(writer, http.StatusBadRequest, map[string]string{"message": "Invalid user id"})
		return
	}

	for i, user := range data.Users {
		if user.ID == id {
			data.Users = append(data.Users[:i], data.Users[i+1:]...)
			responseWithJson(writer, http.StatusOK, map[string]string{"message": "User was deleted"})
			return
		}
	}

	responseWithJson(writer, http.StatusNotFound, map[string]string{"message": "User not found"})
}

func responseWithJson(writer http.ResponseWriter, status int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(object)
}

func generateId(users []dto.User) int {
	var maxId int
	for _, todo := range users {
		if todo.ID > maxId {
			maxId = todo.ID
		}
	}
	return maxId + 1
}
