package mahjong

import (
	"net/http"
)

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
