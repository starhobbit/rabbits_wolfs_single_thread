package movement

import "math/rand"

func FixPos(pos int, max int) int {
	if pos < 0 {
		return 0
	}
	if pos > max-1 {
		return max - 1
	}
	return pos
}

type Position struct {
	X, Y int
}

type MoveDirection struct {
	left, right, top, bottom bool
}

func (d *MoveDirection) Random() {
	d.pairRandom(&d.top, &d.bottom)
	d.pairRandom(&d.left, &d.right)
}

func (d *MoveDirection) pairRandom(first *bool, second *bool) {
	random := rand.Intn(9998)
	switch {
	case random < 3333:
		break
	case random < 6666:
		*first = true
	case random >= 6666:
		*second = true
	}
}

func (p Position) IsValid(max int) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < max && p.Y < max
}

func (p Position) Move(direction MoveDirection, speed uint8, max int) Position {
	if direction.right {
		p.X += int(speed)
	}
	if direction.left {
		p.X -= int(speed)
	}
	if direction.top {
		p.Y -= int(speed)
	}
	if direction.bottom {
		p.Y += int(speed)
	}
	p.X = FixPos(p.X, max)
	p.Y = FixPos(p.Y, max)
	return p
}
