package systems_test

import (
	"testing"

	"github.com/deeean/go-vector/vector2"
	"github.com/stretchr/testify/assert"
	"github.com/yazmeyaa/rpg_game/ecs/components"
	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
	"github.com/yazmeyaa/rpg_game/ecs/systems"
	"github.com/yazmeyaa/rpg_game/ecs/world"
)

func TestMovementSystem(t *testing.T) {
	world := world.NewWorld(2)

	ecsstorage.RegisterComponent(world.Components, components.Position{}, 100, func() *components.Position {
		return &components.Position{}
	})
	ecsstorage.RegisterComponent(world.Components, components.Movement{}, 100, func() *components.Movement {
		return &components.Movement{}
	})

	movSys := systems.NewMovementSystem(world)
	world.Systems.AddSystem(movSys)

	entityID1 := 1
	entityID2 := 2

	posStore, _ := ecsstorage.GetComponentStorage(world.Components, components.Position{})
	movStore, _ := ecsstorage.GetComponentStorage(world.Components, components.Movement{})

	posStore.Add(entityID1, components.Position{X: 0, Y: 0})
	movStore.Add(entityID1, components.Movement{Velocity: vector2.Vector2{X: 1, Y: 1}})

	posStore.Add(entityID2, components.Position{X: 10, Y: 10})
	movStore.Add(entityID2, components.Movement{Velocity: vector2.Vector2{X: -1, Y: -1}})

	world.Systems.Update()

	pos1, _ := posStore.Get(entityID1)
	assert.Equal(t, components.Position{X: 1, Y: 1}, *pos1, "Entity 1 position should be updated correctly")

	pos2, _ := posStore.Get(entityID2)
	assert.Equal(t, components.Position{X: 9, Y: 9}, *pos2, "Entity 2 position should be updated correctly")
}
