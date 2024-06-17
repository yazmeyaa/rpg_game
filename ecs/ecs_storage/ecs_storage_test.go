package ecsstorage_test

import (
	"sync"
	"testing"

	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
)

const max_entities_size int = 1500

type Position struct{ X, Y float32 }
type Health struct{ hp uint32 }

func TestEcsStorage(t *testing.T) {
	manager := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(manager, Position{}, max_entities_size, func() *Position {
		return &Position{}
	})

	store, exist := ecsstorage.GetComponentStorage(manager, Position{})
	if !exist {
		t.Error("Storage is not registered after registration")
	}

	store.Add(1, Position{
		X: 2,
		Y: 4,
	})

	exist = store.Has(1)
	if !exist {
		t.Error("Component is missing after adding to storage")
	}

	pos, _ := store.Get(1)
	if pos.X != 2 || pos.Y != 4 {
		t.Errorf("Wrong values in recieved from storage component. Expected X: %f, Y: %f, Got X: %f, Y: %f", float32(2), float32(4), pos.X, pos.Y)
	}
}

func TestBitmap(t *testing.T) {
	manager := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(manager, Position{}, max_entities_size, func() *Position {
		return &Position{}
	})
	ecsstorage.RegisterComponent(manager, Health{}, max_entities_size, func() *Health {
		return &Health{}
	})

	positionStore, exist := ecsstorage.GetComponentStorage(manager, Position{})
	if !exist {
		t.Error("Storage is not registered after registration")
	}

	healthStore, exist := ecsstorage.GetComponentStorage(manager, Health{})
	if !exist {
		t.Error("Storage is not registered after registration")
	}

	healthStore.Add(1, Health{420})
	positionStore.Add(1, Position{
		X: 2,
		Y: 4,
	})

	healthStore.Add(2, Health{420})
	positionStore.Add(2, Position{
		X: 2,
		Y: 4,
	})

	positionStore.Add(3, Position{
		X: 2,
		Y: 4,
	})

	positionBitmap := positionStore.Bitmap()
	healthBitmap := healthStore.Bitmap()

	if positionBitmap != nil && healthBitmap != nil {
		bitmapq := positionBitmap.Clone(nil)
		bitmapq.And(healthBitmap)
		count := bitmapq.Count()
		if count != 2 {
			t.Errorf("Not valid components count. Expected 2, got %d", count)
		}
	} else {
		t.Error("Bitmap of registered components is not exist")
	}

	if !positionBitmap.Contains(1) {
		t.Error("Bitmap is not contains new component")
	}
}

func TestConcurrentAccess(t *testing.T) {
	const maxEntities = 7800

	manager := ecsstorage.NewComponentsManager()
	ecsstorage.RegisterComponent(manager, Position{}, maxEntities, func() *Position {
		return &Position{}
	})

	store, exist := ecsstorage.GetComponentStorage(manager, Position{})
	if !exist {
		t.Fatal("Storage is not registered after registration")
	}

	var wg sync.WaitGroup

	// Add components concurrently
	for i := 0; i < maxEntities; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			store.Add(id, Position{X: float32(id), Y: float32(id) * 2})
		}(i)
	}

	// Update components concurrently
	for i := 0; i < maxEntities; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			store.Update(id, Position{X: float32(id) * 3, Y: float32(id) * 4})
		}(i)
	}

	// Delete components concurrently
	for i := 0; i < maxEntities; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if id%2 == 0 {
				store.Delete(id)
			}
		}(i)
	}

	wg.Wait()

	// Check final state of the store
	for i := 0; i < maxEntities; i++ {
		comp, exists := store.Get(i)
		if i%2 == 0 {
			if exists {
				t.Errorf("Expected component %d to be deleted", i)
			}
		} else {
			if !exists {
				t.Errorf("Expected component %d to exist", i)
			} else if comp.X != float32(i)*3 || comp.Y != float32(i)*4 {
				t.Errorf("Component %d has wrong value: %v", i, *comp)
			}
		}
	}

	bm := store.Bitmap()

	for i := 0; i < maxEntities; i++ {
		if i%2 == 0 {
			if bm.Contains(uint32(i)) {
				t.Errorf("Bitmap contains ID %d which should be deleted", i)
			}
		} else {
			if !bm.Contains(uint32(i)) {
				t.Errorf("Bitmap does not contain ID %d which should exist", i)
			}
		}
	}
}
