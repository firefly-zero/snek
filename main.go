package main

import (
	"github.com/firefly-zero/firefly-go/firefly"
)

var frame = 0

func init() {
	firefly.Boot = boot
	firefly.Update = update
	firefly.Render = render
}

func boot() {
	apple = NewApple()
	snake = NewSnake()
}

func update() {
	frame += 1
}

func render() {
	firefly.ClearScreen(firefly.ColorWhite)
	apple.Render()
	snake.Render(frame)
}
