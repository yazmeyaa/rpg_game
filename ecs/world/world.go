package world

import (
	"github.com/yazmeyaa/rpg_game/ecs/storage"
)

type World struct {
	Components *storage.ComponentsManager
	Systems    *Systems
}

func NewWorld(systemsCount int) *World {
	return &World{
		Components: storage.NewComponentsManager(),
		Systems:    NewSystems(systemsCount),
	}
}
