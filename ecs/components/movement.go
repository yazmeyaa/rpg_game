package components

import "github.com/deeean/go-vector/vector2"

const (
	MOVEMENT_STORAGE_NAME string = "movement"
)

type Movement struct {
	Velocity     vector2.Vector2
	Acceleration vector2.Vector2
}

func NewDefaultMovement() *Movement {
	return &Movement{}
}
