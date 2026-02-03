package game

import "github.com/firefly-zero/firefly-go/firefly"

const (
	singleplayer firefly.Board = 1
	multiplayer  firefly.Board = 2
)

type Title struct {
	// The text to display.
	msg string

	// If the title screen is blocking, for how much longer it will be displayed.
	// Non-blocking title screen is displayed until it is replaced by a blocking one.
	ttl int

	// Blocking title screen is shown for all peers on the game over.
	// Non-blocking title screen is shown only for one snake when it dies
	// and it's rendered on the background instead of covering the whole screen.
	blocking bool
}

func setTitle(msg string, blocking bool) {
	const defaultTTL = 240

	// If a title is already set, keep it. This way we make sure that if a snake died,
	// we keep the "you died" message instead of  replacing it with "you win" message.
	if title != nil {
		// If the new title is blocking (the "game over" screen)
		// show the blocking title but keep the text.
		if blocking {
			title.ttl = defaultTTL
			title.blocking = true
		}
		return
	}

	if len(snakes.items) == 1 {
		snake := snakes.items[0]
		firefly.AddScore(snake.peer, singleplayer, snake.score.val)
	} else {
		for _, snake := range snakes.items {
			firefly.AddScore(snake.peer, multiplayer, snake.score.val)
		}
	}
	title = &Title{
		msg:      msg,
		ttl:      defaultTTL,
		blocking: blocking,
	}
}

func (t *Title) update() {
	if t.ttl >= 0 {
		t.ttl--
	}
	btns := firefly.ReadButtons(firefly.Combined)
	if t.blocking && (btns.Any() || t.ttl <= 0) {
		resetGame()
	}
}

func (t *Title) render() {
	x := (firefly.Width - font.LineWidth(t.msg)) / 2
	y := (firefly.Height + font.CharHeight()) / 2
	firefly.DrawText(t.msg, font, firefly.P(x, y), firefly.ColorBlack)
}
