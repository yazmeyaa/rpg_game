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
	manager := ecsstorage.NewComponentsManager()
	systems := systems.NewSystems(systemsCount)
	return &World{
		Components: manager,
		Systems:    systems,
	}
}
