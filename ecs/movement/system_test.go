package movement_test

import (
	"testing"

	"github.com/deeean/go-vector/vector2"
	"github.com/stretchr/testify/assert"
	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
	"github.com/yazmeyaa/rpg_game/ecs/movement"
	"github.com/yazmeyaa/rpg_game/ecs/world"
)

func TestMovementSystem(t *testing.T) {
	world := world.NewWorld(2)

	ecsstorage.RegisterComponent(world.Components, movement.Position{}, 100, func() *movement.Position {
		return &movement.Position{}
	})
	ecsstorage.RegisterComponent(world.Components, movement.Movement{}, 100, func() *movement.Movement {
		return &movement.Movement{}
	})

	movSys := movement.NewMovementSystem(world)
	world.Systems.AddSystem(movSys)

	entityID1 := 1
	entityID2 := 2

	posStore, _ := ecsstorage.GetComponentStorage(world.Components, movement.Position{})
	movStore, _ := ecsstorage.GetComponentStorage(world.Components, movement.Movement{})

	posStore.Add(entityID1, movement.Position{X: 0, Y: 0})
	movStore.Add(entityID1, movement.Movement{Velocity: vector2.Vector2{X: 1, Y: 1}})

	posStore.Add(entityID2, movement.Position{X: 10, Y: 10})
	movStore.Add(entityID2, movement.Movement{Velocity: vector2.Vector2{X: -1, Y: -1}})

	world.Systems.Update()

	pos1, _ := posStore.Get(entityID1)
	assert.Equal(t, movement.Position{X: 1, Y: 1}, *pos1, "Entity 1 position should be updated correctly")

	pos2, _ := posStore.Get(entityID2)
	assert.Equal(t, movement.Position{X: 9, Y: 9}, *pos2, "Entity 2 position should be updated correctly")
}
