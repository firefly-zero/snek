package game

import (
	"github.com/firefly-zero/firefly-go/firefly"
)

var (
	font   firefly.Font
	frame  = 0
	snakes *Snakes
)

func Boot() {
	font = firefly.LoadFile("font", nil).Font()
	snakes = newSnakes()
	apple = newApple()
	score = newScore()
}

func Update() {
	frame += 1
	snakes.update()
}

func Render() {
	firefly.ClearScreen(firefly.ColorWhite)
	apple.render()
	snakes.render()
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
