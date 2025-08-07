package game

import (
	"github.com/firefly-zero/firefly-go/firefly"
)

var frame = 0
var font firefly.Font

func Boot() {
	font = firefly.LoadFile("font", nil).Font()
	apple = newApple()
	peers := firefly.GetPeers()
	snakes = make([]*Snake, peers.Len())
	for i, peer := range peers.Slice() {
		snakes[i] = newSnake(peer)
	}
	score = newScore()
}

func Update() {
	frame += 1
	for _, snake := range snakes {
		snake.update(frame, &apple)
		snake.tryEat(&apple, &score)
		score.update(snake)
	}
}

func Render() {
	firefly.ClearScreen(firefly.ColorWhite)
	apple.render()
	for _, snake := range snakes {
		snake.render(frame)
	}
	score.render()
}

func Cheat(c, v int) int {
	switch c {
	case 1:
		apple.move()
		return 1
	case 2:
		for i := 0; i < int(v); i++ {
			score.inc()
		}
		return score.val
	case 3:
		for i := 0; i < int(v); i++ {
			score.dec()
		}
		return score.val
	default:
		return 0
	}
}
