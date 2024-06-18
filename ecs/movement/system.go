package movement

import (
	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
	"github.com/yazmeyaa/rpg_game/ecs/world"
)

type MovementSystem struct {
	positionStorage *ecsstorage.ComponentStorage[Position]
	movementStorage *ecsstorage.ComponentStorage[Movement]
}

func NewMovementSystem(world *world.World) *MovementSystem {
	pStore, _ := ecsstorage.GetComponentStorage(world.Components, Position{})
	mStore, _ := ecsstorage.GetComponentStorage(world.Components, Movement{})
	return &MovementSystem{
		positionStorage: pStore,
		movementStorage: mStore,
	}
}

func (s *MovementSystem) Compute() {
	bitmap := s.positionStorage.Bitmap()
	bitmap.And(s.movementStorage.Bitmap())

	bitmap.Range(func(x uint32) {
		pos, _ := s.positionStorage.Get(int(x))
		mov, _ := s.movementStorage.Get(int(x))
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y
	})
}

func (s *MovementSystem) Priority() uint8 {
	return 1
}
