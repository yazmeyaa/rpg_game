package movement

import "github.com/rpg_game/ecs"

type MovementSystem struct {
}

func NewMovementSystem() ecs.ECSSystem {
	return &MovementSystem{}
}

func (s MovementSystem) Compute() {

}
