package ecsstorage_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
	"github.com/yazmeyaa/rpg_game/ecs/movement"
)

type TestComponent struct {
	Value int
}

func newTestComponent() *TestComponent {
	return &TestComponent{}
}

func TestComponentsManager_Save(t *testing.T) {
	cm := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(cm, TestComponent{}, 10, newTestComponent)
	ecsstorage.RegisterComponent(cm, movement.Position{}, 10, func() *movement.Position {
		return &movement.Position{}
	})

	component := TestComponent{Value: 42}
	storage, _ := ecsstorage.GetComponentStorage(cm, TestComponent{})
	storage.Add(1, component)

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

	newCm := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(newCm, TestComponent{}, 10, newTestComponent)
	ecsstorage.RegisterComponent(newCm, movement.Position{}, 10, func() *movement.Position {
		return &movement.Position{}
	})

	newCm.Load(savePath)
	store, exist := ecsstorage.GetComponentStorage(newCm, TestComponent{})
	assert.True(t, exist, "Store must exist after load")
	pos1, exist := store.Get(1)
	assert.True(t, exist, "Component must exist after load")
	assert.Equal(t, pos1, &TestComponent{Value: 42})
}
