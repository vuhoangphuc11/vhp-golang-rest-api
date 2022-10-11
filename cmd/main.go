package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vuhoangphuc11/vhp-golang-rest-api/pkg/handler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	//vhp: Start Connect Database
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://root:root@cluster0.7wrhdyv.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer client.Disconnect(ctx)

	fmt.Println("Connected to MongoDB!")
	//vhp: End Connect Database

	r := mux.NewRouter()

	r.HandleFunc("/api/get-all-user", handler.GetAllUser).Methods(http.MethodGet)
	r.HandleFunc("/api/get-user/{id}", handler.GetUserById).Methods(http.MethodGet)
	r.HandleFunc("/api/create-user", handler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/api/update-user/{id}", handler.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/api/delete-user/{id}", handler.DeleteUser).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8087", r))

}

func homeLink(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome home!")
}
