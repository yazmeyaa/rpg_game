package ecs

import (
	"reflect"
	"sync"
)

type Id uint16

type Storage struct {
	registry map[string]any
}

func NewStorage() *Storage {
	return &Storage{
		registry: make(map[string]any),
	}
}

func addBasicStorage[T any](s *Storage, dist T) {
	storageName := getName(dist)
	_, exist := s.registry[storageName]
	if !exist {
		s.registry[storageName] = NewBasicStorage(dist)
	}
}

func getName(e interface{}) string {
	typeString := reflect.TypeOf(e).String()
	if typeString[0] == '*' {
		return typeString[1:]
	}
	return typeString
}

func GetStorage[T any](s *Storage, dist T) *BasicStorage[T] {
	storageName := getName(dist)
	regStorage, exist := s.registry[storageName]
	if !exist {
		addBasicStorage(s, dist)
		regStorage = s.registry[storageName]
	}

	storage, ok := regStorage.(*BasicStorage[T])
	if !ok {
		panic("unexpected type in storage registry")
	}
	return storage
}

type BasicStorage[T any] struct {
	items   map[Id]T
	itemsId Id
	mx      sync.RWMutex
}

func NewBasicStorage[T any](dist T) *BasicStorage[T] {
	return &BasicStorage[T]{
		items:   make(map[Id]T),
		itemsId: 0,
		mx:      sync.RWMutex{},
	}
}

func (s *BasicStorage[T]) NewId() Id {
	id := s.itemsId
	s.itemsId++
	return id
}

func (s *BasicStorage[T]) SetItem(id Id, item T) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	s.items[id] = item
}

func (s *BasicStorage[T]) GetItem(id Id) (T, bool) {
	item, exist := s.items[id]
	return item, exist
}

func (s *BasicStorage[T]) AddItem(e T) Id {
	s.mx.RLock()
	defer s.mx.RUnlock()

	itemId := s.NewId()
	s.items[itemId] = e
	return itemId
}

func (s *BasicStorage[T]) GetAllItems() []T {
	list := make([]T, 0, len(s.items))
	for _, item := range s.items {
		list = append(list, item)
	}
	return list
}

func (s *BasicStorage[T]) RemoveItem(id Id) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	delete(s.items, id)
}
