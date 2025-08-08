package game

import "github.com/firefly-zero/firefly-go/firefly"

type Title struct {
	msg string
	ttl int
}

func setTitle(msg string) {
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
	firefly.DrawText(t.msg, font, firefly.P(x, 80), firefly.ColorBlack)
}
