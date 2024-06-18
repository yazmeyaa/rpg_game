package world

import (
	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
)

type World struct {
	Components *ecsstorage.ComponentsManager
	Systems    *Systems
}

func NewWorld(systemsCount int) *World {
	return &World{
		Components: ecsstorage.NewComponentsManager(),
		Systems:    NewSystems(systemsCount),
	}
}
