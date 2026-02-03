package game

import "github.com/firefly-zero/firefly-go/firefly"

const (
	singleplayer firefly.Board = 1
	multiplayer  firefly.Board = 2
)

type Title struct {
	msg string
	ttl int
}

func setTitle(msg string) {
	if len(snakes.items) == 1 {
		snake := snakes.items[0]
		firefly.AddScore(snake.peer, singleplayer, snake.score.val)
	} else {
		for _, snake := range snakes.items {
			firefly.AddScore(snake.peer, multiplayer, snake.score.val)
		}
	}
	title = &Title{
		msg: msg,
		ttl: 240,
	}
}

func (t *Title) update() {
	t.ttl--
	if t.ttl <= 0 {
		resetGame()
	}
}

func (t Title) render() {
	x := (firefly.Width - font.LineWidth(t.msg)) / 2
	y := (firefly.Height + font.CharHeight()) / 2
	firefly.DrawText(t.msg, font, firefly.P(x, y), firefly.ColorBlack)
}
