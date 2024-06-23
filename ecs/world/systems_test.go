package world_test

import (
	"context"
	"testing"
	"time"

	"github.com/deeean/go-vector/vector2"
	"github.com/stretchr/testify/assert"
	"github.com/yazmeyaa/rpg_game/ecs/components"
	"github.com/yazmeyaa/rpg_game/ecs/storage"
	"github.com/yazmeyaa/rpg_game/ecs/systems"
	"github.com/yazmeyaa/rpg_game/ecs/world"
)

func TestWorld(t *testing.T) {
	// ____ SETUP ____
	world := world.NewWorld(1)
	storage.RegisterComponent(world.Components, components.MOVEMENT_STORAGE_NAME, components.Movement{}, 1500, func() *components.Movement {
		return &components.Movement{}
	})
	storage.RegisterComponent(world.Components, components.POSITION_STORAGE_NAME, components.Position{}, 1500, func() *components.Position {
		return &components.Position{}
	})

	movStore, _ := storage.GetComponentStorage[components.Movement](world.Components, components.MOVEMENT_STORAGE_NAME)
	posStore, _ := storage.GetComponentStorage[components.Position](world.Components, components.POSITION_STORAGE_NAME)

	world.Systems.AddSystem(systems.NewMovementSystem(world))
	ctx := context.Background()
	updateTime := time.Duration(time.Second)

	// ____ ADD ENTITY 1 ____

	posStore.Add(1, components.Position{X: 10, Y: 10})
	movStore.Add(1, components.Movement{Velocity: vector2.Vector2{
		X: 1,
		Y: 1,
	}})

	// ____ SIMULATE WORLD ____
	world.Systems.StartUpdating(ctx, updateTime)

	time.Sleep(time.Duration(time.Millisecond * 2500))
	ctx.Done()

	pos, _ := posStore.Get(1)
	assert.True(t, int(pos.X) == 12)
	assert.True(t, int(pos.Y) == 12)
}
