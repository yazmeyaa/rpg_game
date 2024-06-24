package entities

import (
	"github.com/yazmeyaa/rpg_game/ecs/components"
	"github.com/yazmeyaa/rpg_game/ecs/storage"
	"github.com/yazmeyaa/rpg_game/ecs/world"
)

type UserCharacterEntity struct {
	Name     *components.Name
	Position *components.Position
	Movement *components.Movement
}

func NewUserCharacterEntity() *UserCharacterEntity {
	return &UserCharacterEntity{}
}

func (uc *UserCharacterEntity) AddToWorld(world world.World) {

}

func GetUserCharacterByEntityId(entityId int, world *world.World) (UserCharacterEntity, bool) {
	characters, _ := storage.GetComponentStorage[components.UserCharacter](world.Components, components.USERCHARACTER_STORE_NAME)
	if !characters.Has(entityId) {
		return UserCharacterEntity{}, false
	}
	names, _ := storage.GetComponentStorage[components.Name](world.Components, components.NAME_STORAGE_NAME)
	positions, _ := storage.GetComponentStorage[components.Position](world.Components, components.NAME_STORAGE_NAME)
	movements, _ := storage.GetComponentStorage[components.Movement](world.Components, components.NAME_STORAGE_NAME)
	if !names.Has(entityId) || !positions.Has(entityId) || !movements.Has(entityId) {
		return UserCharacterEntity{}, false
	}

	name, _ := names.Get(entityId)
	position, _ := positions.Get(entityId)
	movement, _ := movements.Get(entityId)

	return UserCharacterEntity{
		Name:     name,
		Position: position,
		Movement: movement,
	}, true
}
