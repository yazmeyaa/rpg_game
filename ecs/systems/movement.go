package systems

import (
	"time"

	"github.com/yazmeyaa/rpg_game/ecs/components"
	"github.com/yazmeyaa/rpg_game/ecs/storage"
	"github.com/yazmeyaa/rpg_game/ecs/world"
)

type MovementSystem struct {
	positionStorage *storage.ComponentStorage[components.Position]
	movementStorage *storage.ComponentStorage[components.Movement]
}

func NewMovementSystem(world *world.World) *MovementSystem {
	pStore, _ := storage.GetComponentStorage[components.Position](world.Components, components.POSITION_STORAGE_NAME)
	mStore, _ := storage.GetComponentStorage[components.Movement](world.Components, components.MOVEMENT_STORAGE_NAME)
	return &MovementSystem{
		positionStorage: pStore,
		movementStorage: mStore,
	}
}

func (s *MovementSystem) Compute(dt time.Duration) {
	bitmap := s.positionStorage.Bitmap()
	bitmap.And(s.movementStorage.Bitmap())

	bitmap.Range(func(x uint32) {
		pos, _ := s.positionStorage.Get(int(x))
		mov, _ := s.movementStorage.Get(int(x))
		pos.X += mov.Velocity.X * float64(dt.Milliseconds()) / 1000
		pos.Y += mov.Velocity.Y * float64(dt.Milliseconds()) / 1000
	})
}

func (s *MovementSystem) Priority() uint8 {
	return 1
}
