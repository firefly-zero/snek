package main

import "github.com/firefly-zero/firefly-go/firefly"

var apple Apple

const (
	appleRadius   = 5
	appleDiameter = appleRadius * 2
)

type Apple struct {
	// Coordinates of the apple center
	Pos firefly.Point
}

func NewApple() Apple {
	a := Apple{}
	a.Move()
	return a
}

// move the apple into a new place
func (a *Apple) Move() {
	a.Pos = firefly.Point{
		X: int(firefly.GetRandom()%(firefly.Width-appleRadius*2)) + appleRadius,
		Y: int(firefly.GetRandom()%(firefly.Height-appleRadius*2)) + appleRadius,
	}
}

func (a *Apple) Render() {
	firefly.DrawCircle(
		firefly.Point{X: a.Pos.X - appleRadius, Y: a.Pos.Y - appleRadius},
		appleDiameter,
		firefly.Style{FillColor: firefly.ColorRed},
	)
	firefly.DrawLine(
		a.Pos,
		firefly.Point{X: a.Pos.X + appleRadius, Y: a.Pos.Y - appleRadius},
		firefly.LineStyle{Color: firefly.ColorGreen, Width: 3},
	)
}
