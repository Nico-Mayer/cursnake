package snake

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/internal/geometry"
	"github.com/nico-mayer/cursnake/settings"
	"github.com/nico-mayer/cursnake/sound"
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
		go sound.GetManager().Play("move.mp3")
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
	newHead := head.Add(sb.Direction)
	if *settings.GetSettings().OpenWalls {
		newHead = newHead.Mod(width, height)
	}
	sb.Parts = append(sb.Parts, newHead)
	sb.lastMove -= sb.movementDelay

	if !grow {
		sb.Parts = sb.Parts[1:]
	}
}

func (sb *SnakeBody) Render(screen tcell.Screen) {
	maxGreen := int32(255)
	minGreen := int32(80)

	for i, part := range sb.Parts {
		green := int32(minGreen) + int32(i*2)
		if green > maxGreen {
			green = maxGreen
		}
		screen.SetContent(part.X, part.Y, ' ', nil,
			tcell.StyleDefault.Background(tcell.NewRGBColor(0, green, 0)))
	}
}

func (sb *SnakeBody) CheckSelfCollision() (collided bool) {
	head := sb.GetHead()
	body := sb.Parts[:len(sb.Parts)-1]

	for _, part := range body {
		if part == head {
			return true
		}
	}

	return false
}

func (sb *SnakeBody) CheckWallCollision(screen tcell.Screen) (collided bool) {
	if *settings.GetSettings().OpenWalls {
		return false
	}

	width, height := screen.Size()
	head := sb.GetHead()

	if head.X > width-1 || head.X < 1 || head.Y > height-1 || head.Y < 1 {
		return true
	}

	return false
}

func (sb *SnakeBody) GetHead() geometry.Point {
	return sb.Parts[len(sb.Parts)-1]
}
