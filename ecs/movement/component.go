package movement

import "github.com/deeean/go-vector/vector2"

type Movement struct {
	Velocity     vector2.Vector2
	Acceleration vector2.Vector2
}

type Position struct {
	X float64
	Y float64
}
