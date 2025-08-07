package game

import "github.com/firefly-zero/firefly-go/firefly"

var apple Apple

const (
	appleRadius   = 5
	appleDiameter = appleRadius * 2
)

type Apple struct {
	// Coordinates of the apple center
	pos firefly.Point
}

func newApple() Apple {
	a := Apple{}
	a.move()
	return a
}

// move the apple into a new place
func (a *Apple) move() {
	pos := randomPoint()
	// Don't place the apple inside the snake
	for snakes.appleInside(pos) {
		pos = randomPoint()
	}
	a.pos = pos
}

// Pick a random point for a new apple so that it's fully within the screen.
func randomPoint() firefly.Point {
	x := int(firefly.GetRandom()%(firefly.Width-appleRadius*2)) + appleRadius
	y := int(firefly.GetRandom()%(firefly.Height-appleRadius*2)) + appleRadius
	return firefly.P(x, y)
}

func (a Apple) render() {
	firefly.DrawCircle(
		firefly.Point{X: a.pos.X - appleRadius, Y: a.pos.Y - appleRadius},
		appleDiameter,
		firefly.Solid(firefly.ColorRed),
	)
	firefly.DrawLine(
		a.pos,
		firefly.Point{X: a.pos.X + appleRadius, Y: a.pos.Y - appleRadius},
		firefly.LineStyle{Color: firefly.ColorGreen, Width: 3},
	)
}
