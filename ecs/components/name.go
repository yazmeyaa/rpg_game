package components

const (
	NAME_STORAGE_NAME string = "name"
)

type Name struct {
	Name string
}

func NewDefaultName() *Name {
	return &Name{}
}
