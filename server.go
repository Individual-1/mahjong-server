package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
)

var app *firebase.App
var authClient *auth.Client
var dbClient *db.Client

func main() {
	var err error
	app, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error initializing application: %v\n", err)
	}

	authClient, err = app.Auth()
	if err != nil {
		log.Fatalf("Error initializing auth client: %v\n", err)
	}
	fmt.Println(authClient.VerifyIDToken("asda"))

	dbClient, err = app.Database()
	if err != nil {
		log.Fatalf("Error establishing db connection: %v\n", err)
	}
}

func joinRoom(roomId string) {
	dbClient.NewRef()
}
