package storage

import (
	"sync"

	"github.com/kelindar/bitmap"
)

type ComponentStorer interface {
	Name() string
	ToJSON() ([]byte, error)
	Load(data []byte) error
}

type ComponentsManager struct {
	stores map[string]any
	bitmap map[string]*bitmap.Bitmap
}

func NewComponentsManager() *ComponentsManager {
	return &ComponentsManager{
		stores: map[string]any{},
		bitmap: make(map[string]*bitmap.Bitmap),
	}
}

func (cm *ComponentsManager) IterateOverStores(fn func(key string, store ComponentStorer)) {
	for key, value := range cm.stores {
		if storer, ok := value.(ComponentStorer); ok {
			fn(key, storer)
		} else {
			continue
		}
	}
}

func (cm *ComponentsManager) getBitmap(storeName string) (*bitmap.Bitmap, bool) {
	bm, exist := cm.bitmap[storeName]
	return bm, exist
}

func RegisterComponent[T any](cm *ComponentsManager, storeName string, component T, max_entities_size int, newFunc func() *T) {
	if _, exist := cm.stores[storeName]; exist {
		return
	}

	cm.stores[storeName] = NewComponentStorage(cm, newFunc, storeName)
	bm := &bitmap.Bitmap{}
	bm.Grow(uint32(max_entities_size))
	cm.bitmap[storeName] = bm
}

type ComponentStorageRepository[T any] interface {
	Get(int) (*T, bool)
	Has(entityId int) bool
	Add(entityId int, component T)
	Delete(entityId int)
	Update(entityId int, val T)
	Bitmap() bitmap.Bitmap
	Name() string
}

type ComponentStorage[T any] struct {
	name       string
	components map[int]*T
	pool       sync.Pool
	cm         *ComponentsManager
	mx         sync.RWMutex
}

func NewComponentStorage[T any](cm *ComponentsManager, newFunc func() *T, storeName string) ComponentStorageRepository[T] {
	return &ComponentStorage[T]{
		name:       storeName,
		components: make(map[int]*T),
		pool: sync.Pool{
			New: func() any {
				return newFunc()
			},
		},
		cm: cm,
	}
}

func GetComponentStorage[T any](cm *ComponentsManager, storeName string) (*ComponentStorage[T], bool) {
	storage, exist := cm.stores[storeName]
	if !exist {
		return nil, false
	}

	compStorage, ok := storage.(*ComponentStorage[T])
	if !ok {
		return nil, false
	}

	return compStorage, true
}

func (cs *ComponentStorage[T]) Get(entityId int) (*T, bool) {
	cs.mx.RLock()
	defer cs.mx.RUnlock()

	var element *T
	if !cs.Bitmap().Contains(uint32(entityId)) {
		return element, false
	}

	return cs.components[entityId], true
}

func (cs *ComponentStorage[T]) Has(entityId int) bool {
	cs.mx.RLock()
	defer cs.mx.RUnlock()

	return cs.Bitmap().Contains(uint32(entityId))
}

func (cs *ComponentStorage[T]) Add(entityId int, component T) {
	bm, _ := cs.cm.getBitmap(cs.name)

	cs.mx.Lock()
	defer cs.mx.Unlock()

	if bm.Contains(uint32(entityId)) {
		return
	}

	poolComponent := cs.pool.Get()
	if comp, ok := poolComponent.(*T); ok {
		*comp = component
		cs.components[entityId] = comp
		bm.Set(uint32(entityId))
	}
}

func (cs *ComponentStorage[T]) Delete(entityId int) {
	bm, _ := cs.cm.getBitmap(cs.name)

	cs.mx.Lock()
	defer cs.mx.Unlock()

	if !bm.Contains(uint32(entityId)) {
		return
	}

	component := cs.components[entityId]
	cs.pool.Put(component)
	bm.Remove(uint32(entityId))
	delete(cs.components, entityId)
}

func (cs *ComponentStorage[T]) Update(entityId int, val T) {
	cs.mx.Lock()
	defer cs.mx.Unlock()

	if value, ok := cs.components[entityId]; ok {
		*value = val
	}
}

func (cs *ComponentStorage[T]) Bitmap() bitmap.Bitmap {
	bitmap, _ := cs.cm.getBitmap(cs.name)

	cs.mx.RLock()
	defer cs.mx.RUnlock()

	clone := bitmap.Clone(nil)
	return clone
}

func (cs *ComponentStorage[T]) Name() string {
	return cs.name
}
