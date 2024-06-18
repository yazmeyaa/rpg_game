package ecsstorage_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
	"github.com/yazmeyaa/rpg_game/ecs/movement"
)

func TestSerializeData(t *testing.T) {
	var max_entities_size int = 20
	manager := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(manager, movement.Position{}, max_entities_size, func() *movement.Position {
		return &movement.Position{}
	})

	store, _ := ecsstorage.GetComponentStorage(manager, movement.Position{})

	store.Add(1, movement.Position{X: 2, Y: 2})
	pos, exist := store.Get(1)
	assert.True(t, exist)
	assert.Equal(t, &movement.Position{X: 2, Y: 2}, pos)

	data, err := store.Serialize()
	assert.NoError(t, err, "Ошибка при сериализации данных")

	var deserializedData map[string]map[int]movement.Position

	err = json.Unmarshal(data, &deserializedData)
	assert.NoError(t, err, "Ошибка при десериализации данных")

	assert.Equal(t, movement.Position{X: 2, Y: 2}, deserializedData["components"][1])
}

func TestLoad(t *testing.T) {
	var max_entities_size int = 20
	manager := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(manager, movement.Position{}, max_entities_size, func() *movement.Position {
		return &movement.Position{}
	})

	store, _ := ecsstorage.GetComponentStorage(manager, movement.Position{})

	store.Add(1, movement.Position{X: 2, Y: 2})
	pos, exist := store.Get(1)
	assert.True(t, exist)
	assert.Equal(t, &movement.Position{X: 2, Y: 2}, pos)

	data, err := store.Serialize()
	assert.NoError(t, err, "Ошибка при сериализации данных")

	newManager := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(newManager, movement.Position{}, max_entities_size, func() *movement.Position {
		return &movement.Position{}
	})

	newStore, _ := ecsstorage.GetComponentStorage(newManager, movement.Position{})
	fmt.Println(string(data), newStore)

	laodError := newStore.Load(data)
	assert.NoError(t, laodError, "Erorr while loading error")

	newPos, exist := newStore.Get(1)
	assert.True(t, exist, "Loaded entity is not exist")
	assert.Equal(t, &movement.Position{X: 2, Y: 2}, newPos)
}
