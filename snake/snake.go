package snake

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/internal/geometry"
)

type SnakeBody struct {
	Parts         []geometry.Point
	Direction     geometry.Point
	style         tcell.Style
	lastMove      time.Duration
	movementDelay time.Duration
}

var (
	Up    = geometry.Point{X: 0, Y: -1}
	Down  = geometry.Point{X: 0, Y: 1}
	Left  = geometry.Point{X: -1, Y: 0}
	Right = geometry.Point{X: 1, Y: 0}

	X_SPEED = 35 * time.Millisecond
	Y_SPEED = 60 * time.Millisecond
)

func NewSnakeBody(startX, startY, length int) *SnakeBody {
	parts := make([]geometry.Point, length)
	for i := 0; i < length; i++ {
		parts[i] = geometry.Point{X: startX + i + 1, Y: startY}
	}

	return &SnakeBody{
		Parts:         parts,
		Direction:     Right,
		style:         tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorGreen),
		movementDelay: X_SPEED,
	}
}

func (sb *SnakeBody) setDirection(newDirection geometry.Point) {
	if newDirection.Y != 0 {
		sb.movementDelay = Y_SPEED
	} else {
		sb.movementDelay = X_SPEED
	}

	if newDirection.X != -sb.Direction.X && newDirection.Y != -sb.Direction.Y {
		sb.Direction = newDirection
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

func (sb *SnakeBody) Update(delta time.Duration, width, height int, grow bool) {
	sb.lastMove += delta

	if sb.lastMove <= sb.movementDelay && !grow {
		return
	}

	head := sb.Parts[len(sb.Parts)-1]
	newHead := head.Add(sb.Direction).Mod(width, height)
	sb.Parts = append(sb.Parts, newHead)
	sb.lastMove -= sb.movementDelay

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

func (sb *SnakeBody) CheckFruitCollision(fruitPos geometry.Point) (eaten bool) {
	head := sb.Parts[len(sb.Parts)-1]

	return head == fruitPos
}
