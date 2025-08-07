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
	a.pos = firefly.Point{
		X: int(firefly.GetRandom()%(firefly.Width-appleRadius*2)) + appleRadius,
		Y: int(firefly.GetRandom()%(firefly.Height-appleRadius*2)) + appleRadius,
	}
}

func (a Apple) render() {
	firefly.DrawCircle(
		firefly.Point{X: a.pos.X - appleRadius, Y: a.pos.Y - appleRadius},
		appleDiameter,
		firefly.Style{FillColor: firefly.ColorRed},
	)
	firefly.DrawLine(
		a.pos,
		firefly.Point{X: a.pos.X + appleRadius, Y: a.pos.Y - appleRadius},
		firefly.LineStyle{Color: firefly.ColorGreen, Width: 3},
	)
}
