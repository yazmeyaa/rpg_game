package ecsstorage_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yazmeyaa/rpg_game/ecs/components"
	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
)

func TestSerializeData(t *testing.T) {
	var max_entities_size int = 20
	manager := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(manager, components.Position{}, max_entities_size, func() *components.Position {
		return &components.Position{}
	})

	store, _ := ecsstorage.GetComponentStorage(manager, components.Position{})

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
	manager := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(manager, components.Position{}, max_entities_size, func() *components.Position {
		return &components.Position{}
	})

	store, _ := ecsstorage.GetComponentStorage(manager, components.Position{})

	store.Add(1, components.Position{X: 2, Y: 2})
	pos, exist := store.Get(1)
	assert.True(t, exist)
	assert.Equal(t, &components.Position{X: 2, Y: 2}, pos)

	data, err := store.ToJSON()
	assert.NoError(t, err, "Ошибка при сериализации данных")

	newManager := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(newManager, components.Position{}, max_entities_size, func() *components.Position {
		return &components.Position{}
	})

	newStore, _ := ecsstorage.GetComponentStorage(newManager, components.Position{})
	fmt.Println(string(data), newStore)

	laodError := newStore.Load(data)
	assert.NoError(t, laodError, "Erorr while loading error")

	newPos, exist := newStore.Get(1)
	assert.True(t, exist, "Loaded entity is not exist")
	assert.Equal(t, &components.Position{X: 2, Y: 2}, newPos)
}
