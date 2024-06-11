package ecs_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/rpg_game/ecs"
)

type userCharacterEntity struct {
	name string
	age  uint8
}

type testStorageOne struct{}
type testStorageTwo struct{}
type testStorageThree struct{}

func TestStorage(t *testing.T) {
	store := ecs.NewStorage()
	ecs.SetupStorage(store, userCharacterEntity{}, 10)
	seStorage, _ := ecs.GetStorage(store, userCharacterEntity{})

	seStorage.AddItem(userCharacterEntity{
		name: "Eugene",
		age:  26,
	})

	seStorage.AddItem(userCharacterEntity{
		name: "Oleg",
		age:  60,
	})
}

func TestStorageUpdate(t *testing.T) {
	store := ecs.NewStorage()
	ecs.SetupStorage(store, userCharacterEntity{}, 10)
	seStorage, _ := ecs.GetStorage(store, userCharacterEntity{})

	userId := seStorage.AddItem(userCharacterEntity{
		name: "Eugene",
		age:  24,
	})

	var expectedName string = "Oleg"
	var expectedAge uint8 = 45

	seStorage.SetItem(userId, userCharacterEntity{
		name: expectedName,
		age:  expectedAge,
	})

	user, _ := seStorage.GetItem(userId)
	if user.name != "Oleg" || user.age != 45 {
		t.Errorf("Error while update entity. Expected name=%s, age=%d; recieved name=%s, age=%d", expectedName, expectedAge, user.name, user.age)
	}
}

func TestStorageDelete(t *testing.T) {
	store := ecs.NewStorage()
	ecs.SetupStorage(store, userCharacterEntity{}, 10)
	seStorage, _ := ecs.GetStorage(store, userCharacterEntity{})

	userId := seStorage.AddItem(userCharacterEntity{
		name: "Eugene",
		age:  24,
	})

	seStorage.RemoveItem(userId)
	_, exist := seStorage.GetItem(userId)
	if exist {
		t.Errorf("Entity exist after remove operation")
	}
}

func TestCreateManyStorages(t *testing.T) {
	store := ecs.NewStorage()
	ecs.SetupStorage(store, testStorageOne{}, 10)
	ecs.SetupStorage(store, testStorageTwo{}, 10)
	ecs.SetupStorage(store, testStorageThree{}, 10)
	s1, _ := ecs.GetStorage(store, testStorageOne{})
	s2, _ := ecs.GetStorage(store, testStorageTwo{})
	s3, _ := ecs.GetStorage(store, testStorageThree{})

	s1.AddItem(testStorageOne{})
	s2.AddItem(testStorageTwo{})
	s3.AddItem(testStorageThree{})

	if len(s1.GetAllItems()) != 1 || len(s2.GetAllItems()) != 1 || len(s3.GetAllItems()) != 1 {
		t.Error("Every created storages must have only one entity.")
	}
}

func TestStorageParallelOperations(t *testing.T) {
	storage := ecs.NewStorage()
	ecs.SetupStorage(storage, userCharacterEntity{}, 120)
	userStorage, _ := ecs.GetStorage(storage, userCharacterEntity{})
	wg := sync.WaitGroup{}
	wg.Add(100)

	userId := userStorage.AddItem(userCharacterEntity{
		name: "Eugene",
		age:  42,
	})

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			user, _ := userStorage.GetItem(userId)
			user.name = fmt.Sprintf("%s_%d", user.name, i)
			userStorage.SetItem(userId, user)
		}()
	}

	wg.Wait()
}
