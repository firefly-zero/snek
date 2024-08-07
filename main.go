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
	firefly.Cheat = cheat
}

func boot() {
	font = firefly.LoadROMFile("font").Font()
	apple = NewApple()
	peers := firefly.GetPeers()
	snakes = make([]*Snake, peers.Len())
	for i, peer := range peers.Slice() {
		snakes[i] = NewSnake(peer)
	}
	score = NewScore()
}

func update() {
	frame += 1
	for _, snake := range snakes {
		snake.Update(frame, &apple)
		snake.TryEat(&apple, &score)
		score.Update(snake)
	}
}

func render() {
	firefly.ClearScreen(firefly.ColorWhite)
	apple.Render()
	for _, snake := range snakes {
		snake.Render(frame)
	}
	score.Render()
}

func cheat(c, v int) int {
	switch c {
	case 1:
		apple.Move()
		return 1
	case 2:
		for i := 0; i < int(v); i++ {
			score.Inc()
		}
		return score.val
	case 3:
		for i := 0; i < int(v); i++ {
			score.Dec()
		}
		return score.val
	default:
		return 0
	}
}
