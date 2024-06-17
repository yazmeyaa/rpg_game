package movement

import (
	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
	"github.com/yazmeyaa/rpg_game/ecs/systems"
	"github.com/yazmeyaa/rpg_game/ecs/world"
)

type MovementSystem struct {
	world *world.World
}

func NewMovementSystem(world *world.World) systems.System {
	return &MovementSystem{
		world: world,
	}
}

func (s MovementSystem) Compute() {
	positionStore, _ := ecsstorage.GetComponentStorage(&s.world.Components, Position{})
	movementStore, _ := ecsstorage.GetComponentStorage(&s.world.Components, Movement{})

	bitmap := positionStore.Bitmap()
	bitmap.And(movementStore.Bitmap())

	bitmap.Range(func(x uint32) {
		pos, _ := positionStore.Get(int(x))
		mov, _ := movementStore.Get(int(x))
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y
	})
}

func (s MovementSystem) Priority() uint8 {
	return 1
}
