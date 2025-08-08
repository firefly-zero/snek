package game

import "github.com/firefly-zero/firefly-go/firefly"

type Snakes struct {
	items []*Snake
}

func newSnakes() *Snakes {
	peers := firefly.GetPeers()
	snakes := make([]*Snake, peers.Len())
	for i, peer := range peers.Slice() {
		snakes[i] = newSnake(peer)
	}
	return &Snakes{snakes}
}
func (ss *Snakes) update() {
	if ss == nil {
		return
	}
	for _, snake := range ss.items {
		snake.update()
		snake.tryEat()
	}

	for i, s1 := range snakes.items {
		for j, s2 := range snakes.items {
			sameSnek := i == j
			if s1.bites(sameSnek, s2) {
				if sameSnek {
					firefly.AddProgress(s1.peer, badgeBiteSelf, 1)
				}
				s1.eye.hurt = true
				s1.score.dec()
				if s1.score.val == 0 {
					if sameSnek {
						if s1.peer == me {
							setTitle("u bit urself :(")
						} else {
							setTitle("other snek bit itself, u win")
						}
					} else {
						if s1.peer == me {
							setTitle("u lose :(")
						} else {
							setTitle("u win")
						}
					}
				}
			}
		}
	}

}

func (ss *Snakes) render() {
	if ss == nil {
		return
	}
	for _, snake := range ss.items {
		snake.render()
	}
}

// Check if an apple placed at the given point would collide with any snake.
//
// Used to pick a spot for a new apple position.
func (ss *Snakes) appleInside(pos firefly.Point) bool {
	if ss == nil {
		return false
	}
	for _, s := range ss.items {
		if s.appleCollides(pos) {
			return true
		}
	}
	return false
}
