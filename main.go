package main

import (
	"snek/game"

	"github.com/firefly-zero/firefly-go/firefly"
)

func init() {
	firefly.Boot = game.Boot
	firefly.Update = game.Update
	firefly.Render = game.Render
	firefly.Cheat = game.Cheat
}
