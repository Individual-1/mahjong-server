package defs

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

const (
	East  string = "east"
	South string = "south"
	West  string = "west"
	North string = "north"
)

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
