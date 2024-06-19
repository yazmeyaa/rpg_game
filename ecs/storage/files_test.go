package storage_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yazmeyaa/rpg_game/ecs/components"
	"github.com/yazmeyaa/rpg_game/ecs/storage"
)

type TestComponent struct {
	Value int
}

func newTestComponent() *TestComponent {
	return &TestComponent{}
}

func TestComponentsManager_Save(t *testing.T) {
	cm := storage.NewComponentsManager()
	storage.RegisterComponent(cm, "test_component", TestComponent{}, 10, newTestComponent)
	storage.RegisterComponent(cm, "position", components.Position{}, 10, func() *components.Position {
		return &components.Position{}
	})

	component := TestComponent{Value: 42}
	store, _ := storage.GetComponentStorage[TestComponent](cm, "test_component")
	store.Add(1, component)

	tempDir := t.TempDir()
	savePath := filepath.Join(tempDir, "save.json")

	err := cm.Save(savePath)
	assert.NoError(t, err)

	_, err = os.Stat(savePath)
	assert.NoError(t, err)

	data, err := os.ReadFile(savePath)
	assert.NoError(t, err)

	var loadedData map[string]json.RawMessage
	err = json.Unmarshal(data, &loadedData)
	assert.NoError(t, err)

	newCm := storage.NewComponentsManager()
	storage.RegisterComponent(newCm, "test_component", TestComponent{}, 10, newTestComponent)
	storage.RegisterComponent(newCm, "position", components.Position{}, 10, func() *components.Position {
		return &components.Position{}
	})

	newCm.Load(savePath)
	store, exist := storage.GetComponentStorage[TestComponent](newCm, "test_component")
	assert.True(t, exist, "Store must exist after load")
	pos1, exist := store.Get(1)
	assert.True(t, exist, "Component must exist after load")
	assert.Equal(t, pos1, &TestComponent{Value: 42})
}
