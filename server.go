package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

// User is a struct for a given userID
type User struct {
	Name  string              `firestore:"name"`
	Rooms map[string]UserRoom `firestore:"rooms"`
}

// UserRoom is a struct for roomID within a user object
type UserRoom struct {
	Score int `firestore:"score"`
}

// Room is a struct representing a single game
type Room struct {
	Turn      string            `firestore:"turn"`
	Discard   map[string][]byte `firestore:"discard"`
	Hands     map[string][]byte `firestore:"hands"`
	Pw        string            `firestore:"pw"`
	Users     map[string]string `firestore:"users"`
	Wall      []byte            `firestore:"wall"`
	WallIndex int               `firestore:"wallIndex"`
}

type ErrorType int

const (
	OperationError ErrorType = 1 // Any error from library or other calls
	WashOutError   ErrorType = 5 // Wall is exhausted
)

type MahjongError struct {
	Name string
	Type ErrorType
}

func (e MahjongError) Error() string {
	return e.Name
}

var app *firebase.App
var ctx context.Context
var authClient *auth.Client
var dbClient *firestore.Client

func main() {
	var err error
	ctx = context.Background()
	app, err = firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("Error initializing application: %v\n", err)
	}

	authClient, err = app.Auth()
	if err != nil {
		log.Fatalf("Error initializing auth client: %v\n", err)
	}
	fmt.Println(authClient.VerifyIDToken("asda"))

	// TODO: pull project id from env var or something
	dbClient, err = firestore.NewClient(ctx, "projectID")
	if err != nil {
		log.Fatalf("Error initializing database client: %v\n", err)
	}
}

func joinRoom(roomID string, password string) error {

}

func draw(roomID string, userID string) error {
	var err error
	// TODO: validate userID and roomID

	// Open the room doc
	roomDoc := dbClient.Doc(fmt.Sprintf("rooms/%s", roomID))
	if roomDoc == nil {
		return errors.New("Invalid room ID")
	}

	// Use db transactions to make sure we only update this once
	err = dbClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var room Room
		doc, err := tx.Get(roomDoc)
		if err != nil {
			return MahjongError{Name: err.Error(), Type: OperationError}
		}

		// coerce the doc into our room struct
		if err = doc.DataTo(&room); err != nil {
			return MahjongError{Name: err.Error(), Type: OperationError}
		}

		// Check if user is in the room
		wind, ok := room.Users[userID]
		if !ok {
			return MahjongError{Name: fmt.Sprintf("User %s not in room %s", userID, roomID),
				Type: OperationError}
		}

		if wind != room.Turn {
			return MahjongError{Name: fmt.Sprintf("User %s is not current active player for room %s", userID, roomID),
				Type: OperationError}
		}

		// Check if wall is out of tiles
		if room.WallIndex >= len(room.Wall) {
			return MahjongError{Name: fmt.Sprintf("Room %s is out of Wall tiles", roomID), Type: WashOutError}
		}

		// We meet all pre-reqs, draw a tile and add it to user's hand, then increment
		if room.Hands[wind] == nil {
			room.Hands[wind] = make([]byte, 20)
		}

		// Append tile to user's hand
		room.Hands[wind] = append(room.Hands[wind], room.Wall[room.WallIndex])

		tx.Update(roomDoc, []firestore.Update{
			{Path: "wallIndex", Value: room.WallIndex + 1},
			{Path: "hands", Value: room.Hands},
		})
	})
}

func discard(roomID string, userID string, tileID int) error {
	var err error
	// TODO: validate userID and roomID

	// Open the room doc
	roomDoc := dbClient.Doc(fmt.Sprintf("rooms/%s", roomID))
	if roomDoc == nil {
		return errors.New("Invalid room ID")
	}

	// Use db transactions to make sure we only update this once
	err = dbClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var room Room
		doc, err := tx.Get(roomDoc)
		if err != nil {
			return MahjongError{Name: err.Error(), Type: OperationError}
		}

		// coerce the doc into our room struct
		if err = doc.DataTo(&room); err != nil {
			return MahjongError{Name: err.Error(), Type: OperationError}
		}

		// Check if user is in the room
		wind, ok := room.Users[userID]
		if !ok {
			return MahjongError{Name: fmt.Sprintf("User %s not in room %s", userID, roomID),
				Type: OperationError}
		}

		if wind != room.Turn {
			return MahjongError{Name: fmt.Sprintf("User %s is not current active player for room %s", userID, roomID),
				Type: OperationError}
		}

		// Check if user has tile in their hand
		if !bytes.Contains(room.Hands[userID], []byte(tileID)) {
			return MahjongError{Name: fmt.Sprintf("Room %s is out of Wall tiles", roomID), Type: WashOutError}
		}

		// We meet all pre-reqs, draw a tile and add it to user's hand, then increment
		if room.Hands[wind] == nil {
			room.Hands[wind] = make([]byte, 20)
		}

		// Append tile to user's hand
		room.Hands[wind] = append(room.Hands[wind], room.Wall[room.WallIndex])

		tx.Update(roomDoc, []firestore.Update{
			{Path: "wallIndex", Value: room.WallIndex + 1},
			{Path: "hands", Value: room.Hands},
		})
	})
}
