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

	. "github.com/Individual-1/mahjong-server/defs"
)

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
			{Path: fmt.Sprintf("hands.%s", wind), Value: room.Hands[wind]},
		})
	})
}

func discard(roomID string, userID string, tileID byte) error {
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

		// Get the index of the tile we want to discard
		index := bytes.IndexByte(room.Hands[userID], tileID)

		if index == -1 {
			return MahjongError{Name: fmt.Sprintf("User %s in room %s does not have tile %i", userID, roomID, tileID),
				Type: OperationError}
		}

		// Remove tile from user's hand
		if index == 0 {
			room.Hands[wind] = room.Hands[wind][1:]
		} else if index != len(room.Hands[wind])-1 {
			room.Hands[wind] = append(room.Hands[wind][:index], room.Hands[wind][index+1:]...)
		} else {
			room.Hands[wind] = room.Hands[wind][:len(room.Hands[wind])-1]
		}

		room.Discard[wind] = append(room.Discard[wind], tileID)

		tx.Update(roomDoc, []firestore.Update{
			{Path: fmt.Sprintf("discard.%s", wind), Value: room.Discard[wind]},
			{Path: fmt.Sprintf("hands.%s", wind), Value: room.Hands[wind]},
		})
	})
}

// Claim a tile, puts your name in to pool and the person who claimed first with the highest priority gets it
func claim(roomID string, userID string, tileID byte) error {

}
