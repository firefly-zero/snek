package game

import (
	"github.com/firefly-zero/firefly-go/firefly"
)

var (
	font   firefly.Font
	frame  int
	snakes *Snakes
	title  *Title
	me     firefly.Peer
)

func Boot() {
	font = firefly.LoadFile("font", nil).Font()
	me = firefly.GetMe()
	resetGame()
}

func resetGame() {
	snakes = newSnakes()
	apple = newApple()
	frame = 0
	title = nil
}

func Update() {
	if title != nil {
		title.update()
		return
	}
	frame += 1
	snakes.update()
}

func Render() {
	firefly.ClearScreen(firefly.ColorWhite)
	if title != nil {
		title.render()
		return
	}
	apple.render()
	snakes.render()
}

func Cheat(c, v int) int {
	switch c {
	case 1:
		apple.move()
		return 1
	case 2:
		s := getMySnake()
		for i := 0; i < int(v); i++ {
			s.score.inc()
		}
		return int(s.score.val)
	case 3:
		s := getMySnake()
		for i := 0; i < int(v); i++ {
			s.score.dec()
		}
		return int(s.score.val)
	default:
		return 0
	}
}

func getMySnake() *Snake {
	for _, s := range snakes.items {
		if s.peer == me {
			return s
		}
	}
	return snakes.items[0]
}
