package world

import (
	"github.com/yazmeyaa/rpg_game/ecs/storage"
)

type World struct {
	Components *storage.ComponentsManager
	Systems    *Systems
}

func NewWorld() *World {
	return &World{
		Components: storage.NewComponentsManager(),
		Systems:    NewSystems(),
	}
}
