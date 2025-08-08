package game

import (
	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

type Eye struct {
	// The point the snake is looking at.
	lookingAt firefly.Point

	// If true, the snake has bumped into another snake or itself on this update.
	// Used to highlight the snake's eye with red.
	hurt bool

	// The timer for the snake's eye blinking.
	blinkCounter int
	blinkMaxTime int
}

func (s *Eye) update(mouth firefly.Point) {
	// Calculate position of eye based on the where the apple is
	lookX := float32(apple.pos.X - mouth.X)
	lookY := float32(apple.pos.Y - mouth.Y)
	lookLen := tinymath.Hypot(lookX, lookY)
	dX := lookX * 3 / lookLen
	dY := lookY * 3 / lookLen

	s.lookingAt = firefly.Point{
		X: mouth.X + int(dX),
		Y: mouth.Y + int(dY),
	}

	s.blinkCounter += int(firefly.GetRandom() % 5)
	if s.blinkCounter > s.blinkMaxTime {
		s.blinkCounter = 0
		s.blinkMaxTime = int(100 + firefly.GetRandom()%100)
	}
}

func (s *Eye) render(mouth firefly.Point) {
	style := firefly.Solid(firefly.ColorWhite)
	if s.hurt {
		style.FillColor = firefly.ColorRed
	}
	// We reset it only after rendering to make sure to render it
	// for at least one frame.
	s.hurt = false

	// Outer dark circle representing the head.
	firefly.DrawCircle(
		firefly.Point{
			X: mouth.X - snakeWidth/2 - 1,
			Y: mouth.Y - snakeWidth/2 - 1,
		},
		snakeWidth+2,
		firefly.Solid(firefly.ColorBlue),
	)

	// Inner light circle representing the open eyelids.
	firefly.DrawCircle(
		firefly.Point{
			X: mouth.X - snakeWidth/2,
			Y: mouth.Y - snakeWidth/2,
		},
		snakeWidth,
		firefly.Solid(firefly.ColorLightBlue),
	)

	// White circle representing the eyeball.
	firefly.DrawCircle(
		firefly.Point{
			X: mouth.X - snakeWidth/2 + 1,
			Y: mouth.Y - snakeWidth/2 + 1,
		},
		snakeWidth-2,
		style,
	)

	// Black circle representing the eye iris.
	firefly.DrawCircle(
		firefly.P(
			s.lookingAt.X-snakeWidth/8,
			s.lookingAt.Y-snakeWidth/8,
		),
		snakeWidth/4,
		firefly.Solid(firefly.ColorBlack),
	)

	// If it's soon time to blink, close the eyelid.
	if s.blinkCounter < 20 {
		firefly.DrawCircle(
			firefly.P(
				mouth.X-snakeWidth/2+1,
				mouth.Y-snakeWidth/2+1,
			),
			snakeWidth-2,
			firefly.Solid(firefly.ColorLightBlue),
		)
	}

}
