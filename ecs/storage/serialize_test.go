package storage_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yazmeyaa/rpg_game/ecs/components"
	"github.com/yazmeyaa/rpg_game/ecs/storage"
)

func TestSerializeData(t *testing.T) {
	var max_entities_size int = 20
	manager := storage.NewComponentsManager()
	storage.RegisterComponent(manager, "position", components.Position{}, max_entities_size, func() *components.Position {
		return &components.Position{}
	})

	store, _ := storage.GetComponentStorage[components.Position](manager, "position")

	store.Add(1, components.Position{X: 2, Y: 2})
	pos, exist := store.Get(1)
	assert.True(t, exist)
	assert.Equal(t, &components.Position{X: 2, Y: 2}, pos)

	data, err := store.ToJSON()
	assert.NoError(t, err, "Ошибка при сериализации данных")

	var deserializedData map[string]map[int]components.Position

	err = json.Unmarshal(data, &deserializedData)
	assert.NoError(t, err, "Ошибка при десериализации данных")

	assert.Equal(t, components.Position{X: 2, Y: 2}, deserializedData["components"][1])
}

func TestLoad(t *testing.T) {
	var max_entities_size int = 20
	manager := storage.NewComponentsManager()
	storage.RegisterComponent(manager, "position", components.Position{}, max_entities_size, func() *components.Position {
		return &components.Position{}
	})

	store, _ := storage.GetComponentStorage[components.Position](manager, "position")

	store.Add(1, components.Position{X: 2, Y: 2})
	pos, exist := store.Get(1)
	assert.True(t, exist)
	assert.Equal(t, &components.Position{X: 2, Y: 2}, pos)

	data, err := store.ToJSON()
	assert.NoError(t, err, "Ошибка при сериализации данных")

	newManager := storage.NewComponentsManager()
	storage.RegisterComponent(newManager, "position", components.Position{}, max_entities_size, func() *components.Position {
		return &components.Position{}
	})

	newStore, _ := storage.GetComponentStorage[components.Position](newManager, "position")

	laodError := newStore.Load(data)
	assert.NoError(t, laodError, "Erorr while loading error")

	newPos, exist := newStore.Get(1)
	assert.True(t, exist, "Loaded entity is not exist")
	assert.Equal(t, &components.Position{X: 2, Y: 2}, newPos)
}
