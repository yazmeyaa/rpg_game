package world

import (
	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
	"github.com/yazmeyaa/rpg_game/ecs/systems"
)

type World struct {
	Components *ecsstorage.ComponentsManager
	Systems    *systems.Systems
}

func NewWorld(systemsCount int) *World {
	return &World{
		Components: ecsstorage.NewComponentsManager(),
		Systems:    systems.NewSystems(systemsCount),
	}
}
