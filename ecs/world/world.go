package world

import (
	ecsstorage "github.com/yazmeyaa/rpg_game/ecs/ecs_storage"
	"github.com/yazmeyaa/rpg_game/ecs/systems"
)

type World struct {
	Components ecsstorage.ComponentsManager
	Systems    systems.Systems
}
