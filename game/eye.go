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

func (eye *Eye) update(mouth firefly.Point) {
	// Calculate position of eye based on the where the apple is
	lookX := float32(apple.pos.X - mouth.X)
	lookY := float32(apple.pos.Y - mouth.Y)
	lookLen := tinymath.Hypot(lookX, lookY)
	dX := lookX * 3 / lookLen
	dY := lookY * 3 / lookLen

	eye.lookingAt = firefly.Point{
		X: mouth.X + int(dX),
		Y: mouth.Y + int(dY),
	}

	eye.blinkCounter += int(firefly.GetRandom() % 5)
	if eye.blinkCounter > eye.blinkMaxTime {
		eye.blinkCounter = 0
		eye.blinkMaxTime = int(100 + firefly.GetRandom()%100)
	}
}

func (eye *Eye) render(mouth firefly.Point, me bool) {
	style := firefly.Solid(firefly.ColorWhite)
	if eye.hurt {
		style.FillColor = firefly.ColorRed
	}
	// We reset it only after rendering to make sure to render it
	// for at least one frame.
	eye.hurt = false

	// Outer dark circle representing the head.
	headColor := firefly.ColorBlue
	if !me {
		headColor = firefly.ColorGray
	}
	firefly.DrawCircle(
		firefly.Point{
			X: mouth.X - snakeWidth/2 - 1,
			Y: mouth.Y - snakeWidth/2 - 1,
		},
		snakeWidth+2,
		firefly.Solid(headColor),
	)

	// Inner light circle representing the open eyelids.
	eyelidColor := firefly.ColorLightBlue
	if !me {
		eyelidColor = firefly.ColorLightGray
	}
	firefly.DrawCircle(
		firefly.Point{
			X: mouth.X - snakeWidth/2,
			Y: mouth.Y - snakeWidth/2,
		},
		snakeWidth,
		firefly.Solid(eyelidColor),
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
			eye.lookingAt.X-snakeWidth/8,
			eye.lookingAt.Y-snakeWidth/8,
		),
		snakeWidth/4,
		firefly.Solid(firefly.ColorBlack),
	)

	// If it's soon time to blink, close the eyelid.
	if eye.blinkCounter < 20 {
		firefly.DrawCircle(
			firefly.P(
				mouth.X-snakeWidth/2+1,
				mouth.Y-snakeWidth/2+1,
			),
			snakeWidth-2,
			firefly.Solid(eyelidColor),
		)
	}
}
