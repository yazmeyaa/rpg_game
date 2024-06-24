package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yazmeyaa/rpg_game/ecs/components"
	"github.com/yazmeyaa/rpg_game/ecs/storage"
	"github.com/yazmeyaa/rpg_game/ecs/systems"
	"github.com/yazmeyaa/rpg_game/ecs/world"
)

func registerWorldComponents(world *world.World, max_entities_size int) {
	storage.RegisterComponent(world.Components, components.MOVEMENT_STORAGE_NAME, components.Movement{}, max_entities_size, components.NewDefaultMovement)
	storage.RegisterComponent(world.Components, components.POSITION_STORAGE_NAME, components.Position{}, max_entities_size, components.NewDefaultPosition)
	storage.RegisterComponent(world.Components, components.NAME_STORAGE_NAME, components.Name{}, max_entities_size, components.NewDefaultName)
	storage.RegisterComponent(world.Components, components.USERCHARACTER_STORE_NAME, components.UserCharacter{}, max_entities_size, func() *components.UserCharacter {
		return &components.UserCharacter{}
	})
}

func registerWorldSystems(world *world.World) {
	world.Systems.AddSystem(systems.NewMovementSystem(world))
}

func setupWorld(world *world.World, max_entities_size int) {
	registerWorldComponents(world, max_entities_size)
	registerWorldSystems(world)
}

func main() {
	world := world.NewWorld()
	const max_entities_size int = 3000
	setupWorld(world, max_entities_size)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updateInterval := time.Millisecond * 16
	world.Systems.StartUpdating(ctx, updateInterval)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	cancel()
	fmt.Fprintln(os.Stdout, []any{"Server gracefully stopped."}...)
}
