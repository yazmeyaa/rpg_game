package components

const (
	POSITION_STORAGE_NAME string = "position"
)

type Position struct {
	X float64
	Y float64
}

func NewDefaultPosition() *Position {
	return &Position{}
}
