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
		snake.update(&apple)
		snake.tryEat(&apple, &score)
		score.update(snake)
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
		if s.collides(pos) {
			return true
		}
	}
	return false
}
