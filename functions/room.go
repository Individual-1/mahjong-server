package mahjong

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
)

// User is a struct for a given userID
type User struct {
	name  string              // name
	rooms map[string]UserRoom // $roomId: UserRoom
}

// UserRoom is a struct for roomID within a user object
type UserRoom struct {
	score int    // score
	tiles []byte // tiles
	wind  int    // wind
}

// Room is a struct representing a single game
type Room struct {
	discard   map[int][]byte  // discard
	pw        string          // pw
	users     map[string]bool // users
	wall      []byte          // wall
	wallIndex int             //wallIndex
}

var app *firebase.App
var authClient *auth.Client
var dbClient *db.Client

func HandleRoom(w http.ResponseWriter, r *http.Request) {
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

// Create room with optional password and return room code or error
func createRoom(password string) (string, error) {
	var err error

	roomsRef := dbClient.NewRef("rooms")
	// TODO: figure out how to generate a non-conflicting room ID without synchronization

}

// Leave specified room
func leaveRoom(roomID string) error {
	var err error

	// TODO: validate roomID string, strict character set and length

	/*
		1. Read user rooms to see if they are a member of this room, if so pull the wind
		2.
	*/

	roomRef := dbClient.NewRef(fmt.Sprintf("rooms/%s", roomID))
}

// Join specified room with optional password
func joinRoom(roomID string, password string) error {

}
