package main

import (
	"github.com/firefly-zero/firefly-go/firefly"
)

var frame = 0
var font firefly.Font

func init() {
	firefly.Boot = boot
	firefly.Update = update
	firefly.Render = render
}

func boot() {
	font = firefly.LoadROMFile("font").Font()
	apple = NewApple()
	snake = NewSnake()
	score = NewScore()
}

func update() {
	frame += 1
	snake.Update(frame, &apple)
	snake.TryEat(&apple, &score)
	score.Update(snake)
}

func render() {
	firefly.ClearScreen(firefly.ColorWhite)
	apple.Render()
	snake.Render(frame)
	score.Render()
}
