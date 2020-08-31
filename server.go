package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
)

// User is a struct for a given userID
type User struct {
	Name  string              `json:"name"`
	Rooms map[string]UserRoom `json:"rooms"`
}

// UserRoom is a struct for roomID within a user object
type UserRoom struct {
	Score int `json:"score"`
}

// Room is a struct representing a single game
type Room struct {
	Turn      string            `json:"turn"`
	Discard   RoomTiles         `json:"discard"`
	Hands     RoomTiles         `json:"hands"`
	Pw        string            `json:"pw"`
	Users     map[string]string `json:"users"`
	Wall      []byte            `json:"wall"`
	WallIndex int               `json:"wallIndex"`
}

type RoomTiles struct {
	North []byte `json:"north"`
	West  []byte `json:"west"`
	South []byte `json:"south"`
	East  []byte `json:"east"`
}

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

}

func joinRoom(roomID string) {

}

func draw(roomID string, userID string) {
	// Use db transactions to make sure we only update this once
}
