package main

import (
	"github.com/firefly-zero/firefly-go/firefly"
)

func init() {
	firefly.Boot = boot
	firefly.Update = update
	firefly.Render = render
}

func boot() {
	apple = NewApple()
}

func update() {
	//...
}

func render() {
	firefly.ClearScreen(firefly.ColorWhite)
	apple.Render()
}
