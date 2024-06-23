package storage

import (
	"encoding/json"
	"fmt"
)

type SerializedData[T any] struct {
	Components map[int]T `json:"components"`
}

func (cs *ComponentStorage[T]) ToJSON() ([]byte, error) {
	a := SerializedData[*T]{
		Components: cs.components,
	}

	outString, error := json.Marshal(a)
	if error != nil {
		return make([]byte, 0), error
	}

	return outString, nil
}

func (cs *ComponentStorage[T]) load(data []byte) error {
	var deserializeData SerializedData[T]
	err := json.Unmarshal(data, &deserializeData)
	if err != nil {
		return fmt.Errorf("ошибка десериализации JSON: %v", err)
	}

	for entityId, component := range deserializeData.Components {
		cs.Add(entityId, component)
	}

	return nil
}

func (cs *ComponentStorage[T]) Load(data []byte) error {
	return cs.load(data)
}
