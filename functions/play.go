package mahjong

import (
	"net/http"
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

func HandleRoom(w http.ResponseWriter, r *http.Request) {

}

// Create room with optional password and return room code or error
func createRoom(password string) (string, error) {

}

// Leave specified room
func leaveRoom(roomID string) error {

}

// Join specified room with optional password
func joinRoom(roomID string, password string) error {

}
