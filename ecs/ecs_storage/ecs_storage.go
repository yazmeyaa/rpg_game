package ecsstorage

import (
	"reflect"
	"sync"

	"github.com/kelindar/bitmap"
	"github.com/yazmeyaa/sparse_set"
)

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

func (cm *ComponentsManager) getBitmap(dist any) (*bitmap.Bitmap, bool) {
	bm, exist := cm.bitmap[getName(dist)]
	return bm, exist
}

func getName(dist any) string {
	return reflect.TypeOf(dist).String()
}

func RegisterComponent[T any](cm *ComponentsManager, component T, max_entities_size int, newFunc func() *T) {
	name := getName(component)
	if _, exist := cm.stores[name]; exist {
		return
	}

	cm.stores[name] = NewComponentStorage(cm, max_entities_size, newFunc)
	bm := &bitmap.Bitmap{}
	bm.Grow(uint32(max_entities_size))
	cm.bitmap[name] = bm
}

type ComponentStorage[T any] struct {
	componentReference T
	name               string
	components         map[int]*T
	sparseSet          sparse_set.SparseSet
	pool               sync.Pool
	cm                 *ComponentsManager
	mx                 sync.RWMutex
}

func NewComponentStorage[T any](cm *ComponentsManager, max_entities_size int, newFunc func() *T) *ComponentStorage[T] {
	return &ComponentStorage[T]{
		componentReference: *newFunc(),
		name:               getName(newFunc()),
		components:         make(map[int]*T),
		sparseSet:          *sparse_set.NewSparseSet(uint32(max_entities_size)),
		pool: sync.Pool{
			New: func() any {
				return newFunc()
			},
		},
		cm: cm,
	}
}

func GetComponentStorage[T any](cm *ComponentsManager, dist T) (*ComponentStorage[T], bool) {
	name := getName(dist)
	storage, exist := cm.stores[name]
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
	if !cs.sparseSet.Contains(entityId) {
		return element, false
	}

	return cs.components[entityId], true
}

func (cs *ComponentStorage[T]) Has(entityId int) bool {
	cs.mx.RLock()
	defer cs.mx.RUnlock()

	return cs.sparseSet.Contains(entityId)
}

func (cs *ComponentStorage[T]) Add(entityId int, component T) {
	cs.mx.Lock()
	defer cs.mx.Unlock()

	if cs.sparseSet.Contains(entityId) {
		return
	}

	poolComponent := cs.pool.Get()
	if comp, ok := poolComponent.(*T); ok {
		*comp = component
		cs.components[entityId] = comp
		cs.sparseSet.Add(entityId)

		bm, _ := cs.cm.getBitmap(component)
		bm.Set(uint32(entityId))
	}
}

func (cs *ComponentStorage[T]) Delete(entityId int) {
	cs.mx.Lock()
	defer cs.mx.Unlock()

	if !cs.sparseSet.Contains(entityId) {
		return
	}

	component := cs.components[entityId]
	cs.pool.Put(component)
	cs.sparseSet.Remove(entityId)
	delete(cs.components, entityId)

	bm, _ := cs.cm.getBitmap(cs.componentReference)
	bm.Remove(uint32(entityId))
}

func (cs *ComponentStorage[T]) Update(entityId int, val T) {
	cs.mx.Lock()
	defer cs.mx.Unlock()

	if value, ok := cs.components[entityId]; ok {
		*value = val
	}
}

func (cs *ComponentStorage[T]) Bitmap() bitmap.Bitmap {
	cs.mx.RLock()
	defer cs.mx.RUnlock()

	bitmap, _ := cs.cm.getBitmap(cs.componentReference)
	clone := bitmap.Clone(nil)
	return clone
}
