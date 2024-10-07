package snake

import (
	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/geometry"
)

type SnakeBody struct {
	Parts     []geometry.Point
	Direction geometry.Point
	style     tcell.Style
}

var (
	Up    = geometry.Point{X: 0, Y: -1}
	Down  = geometry.Point{X: 0, Y: 1}
	Left  = geometry.Point{X: -1, Y: 0}
	Right = geometry.Point{X: 1, Y: 0}
)

func NewSnakeBody(startX, startY, length int) *SnakeBody {
	parts := make([]geometry.Point, length)
	for i := 0; i < length; i++ {
		parts[i] = geometry.Point{X: startX + i + 1, Y: startY}
	}

	return &SnakeBody{
		Parts:     parts,
		Direction: Right,
		style:     tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorGreen),
	}
}

func (sb *SnakeBody) Up() {
	sb.setDirection(Up)
}
func (sb *SnakeBody) Down() {
	sb.setDirection(Down)
}
func (sb *SnakeBody) Left() {
	sb.setDirection(Left)
}
func (sb *SnakeBody) Right() {
	sb.setDirection(Right)
}

func (sb *SnakeBody) setDirection(newDirection geometry.Point) {
	if newDirection.X != -sb.Direction.X && newDirection.Y != -sb.Direction.Y {
		sb.Direction = newDirection
	}
}

func (sb *SnakeBody) Update(width, height int, grow bool) {
	head := sb.Parts[len(sb.Parts)-1]
	newHead := head.Add(sb.Direction).Mod(width, height)
	sb.Parts = append(sb.Parts, newHead)

	if !grow {
		sb.Parts = sb.Parts[1:]
	}
}

func (sb *SnakeBody) Render(screen tcell.Screen) {
	for _, part := range sb.Parts {
		screen.SetContent(part.X, part.Y, ' ', nil, sb.style)
	}
}

func (sb *SnakeBody) CheckSelfCollision() (collided bool) {
	head := sb.Parts[len(sb.Parts)-1]
	body := sb.Parts[:len(sb.Parts)-1]

	for _, part := range body {
		if part == head {
			return true
		}
	}

	return false
}
