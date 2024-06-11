package ecs

import (
	"reflect"
	"sync"

	"github.com/yazmeyaa/sparse_set"
)

type Id uint16

type Storage struct {
	registry map[string]any
	mx       sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		registry: make(map[string]any),
	}
}

func getName(e interface{}) string {
	typeString := reflect.TypeOf(e).String()
	if typeString[0] == '*' {
		return typeString[1:]
	}
	return typeString
}

func GetStorage[T any](s *Storage, dist T) (*BasicStorage[T], bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	storageName := getName(dist)
	regStorage, exist := s.registry[storageName]
	if !exist {
		return nil, false
	}

	storage, ok := regStorage.(*BasicStorage[T])
	if !ok {
		panic("unexpected type in storage registry")
	}
	return storage, true
}
func SetupStorage[T any](s *Storage, dist T, maxSize int) {
	s.mx.Lock()
	defer s.mx.Unlock()

	storageName := getName(dist)
	_, exist := s.registry[storageName]
	if !exist {
		s.registry[storageName] = NewBasicStorage(dist, maxSize)
	}
}

type BasicStorage[T any] struct {
	items     map[Id]T
	sparseSet *sparse_set.SparseSet
	itemsId   Id
	mx        sync.RWMutex
}

func NewBasicStorage[T any](dist T, maxSize int) *BasicStorage[T] {
	return &BasicStorage[T]{
		items:     make(map[Id]T),
		sparseSet: sparse_set.NewSparseSet(uint32(maxSize)),
		itemsId:   0,
		mx:        sync.RWMutex{},
	}
}

func (s *BasicStorage[T]) NewId() Id {
	id := s.itemsId
	s.itemsId++
	return id
}

func (s *BasicStorage[T]) SetItem(id Id, item T) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.items[id] = item
	s.sparseSet.Add(int(id))
}

func (s *BasicStorage[T]) GetItem(id Id) (T, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	item, exist := s.items[id]
	return item, exist
}

func (s *BasicStorage[T]) AddItem(e T) Id {
	s.mx.Lock()
	defer s.mx.Unlock()

	itemId := s.NewId()
	s.items[itemId] = e
	s.sparseSet.Add(int(itemId))
	return itemId
}

func (s *BasicStorage[T]) GetAllItems() []T {
	s.mx.RLock()
	defer s.mx.RUnlock()

	ids := s.sparseSet.GetAll()
	list := make([]T, 0, len(ids))
	for _, id := range s.sparseSet.GetAll() {
		list = append(list, s.items[Id(id)])
	}
	return list
}

func (s *BasicStorage[T]) RemoveItem(id Id) {
	s.mx.Lock()
	defer s.mx.Unlock()

	delete(s.items, id)
	s.sparseSet.Remove(int(id))
}
