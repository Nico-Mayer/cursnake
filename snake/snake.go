package snake

import (
	"math"
	"slices"
	"time"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/internal/geometry"
	"github.com/nico-mayer/cursnake/internal/utils"
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
		go sound.GetManager().Play("click.mp3")
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
	if settings.Get().Wraparound {
		innerWidth := width - 2
		innerHeight := height - 2
		newHead = newHead.Sub(geometry.Point{X: 1, Y: 1}).Mod(innerWidth, innerHeight).Add(geometry.Point{X: 1, Y: 1})
	}
	sb.Parts = append(sb.Parts, newHead)
	sb.lastMove -= sb.movementDelay

	if !grow {
		sb.Parts = sb.Parts[1:]
	}
}

func (sb *SnakeBody) Render(screen tcell.Screen) {
	background, err := utils.HexColorToInt32(settings.Get().SnakeBodyOptions.Background)
	if err != nil {
		background, _ = utils.HexColorToInt32(settings.GetDefaultSettings().SnakeBodyOptions.Background)
	}
	foreground, err := utils.HexColorToInt32(settings.Get().SnakeBodyOptions.Foreground)
	if err != nil {
		foreground, _ = utils.HexColorToInt32(settings.GetDefaultSettings().SnakeBodyOptions.Foreground)
	}

	totalParts := len(sb.Parts)

	for i, part := range sb.Parts {
		reversedIndex := int32(totalParts - 1 - i)
		renderSymbol := settings.Get().SnakeBodyOptions.BodyRune
		if i == totalParts-1 {
			renderSymbol = settings.Get().SnakeBodyOptions.HeadRune
		}
		if len(renderSymbol) == 0 {
			renderSymbol = " "
		}

		r, g, b := tcell.NewHexColor(background).RGB()

		r = int32(utils.Clamp(int(math.Max(0, float64(r-150))), int(r-(reversedIndex*2)), int(r)))
		g = int32(utils.Clamp(int(math.Max(0, float64(g-150))), int(g-(reversedIndex*2)), int(g)))
		b = int32(utils.Clamp(int(math.Max(0, float64(b-150))), int(b-(reversedIndex*2)), int(b)))

		screen.SetContent(part.X, part.Y, rune(renderSymbol[0]), nil,
			tcell.StyleDefault.Background(tcell.NewRGBColor(r, g, b)).Foreground(tcell.NewHexColor(foreground)))
	}
}

func (sb *SnakeBody) CheckSelfCollision() (collided bool) {
	head := sb.GetHead()
	body := sb.Parts[:len(sb.Parts)-1]

	if slices.Contains(body, head) {
		return true
	}

	return false
}

func (sb *SnakeBody) CheckWallCollision(screen tcell.Screen) (collided bool) {
	if settings.Get().Wraparound {
		return false
	}
	offset := 1
	width, height := screen.Size()
	head := sb.GetHead()

	if head.X >= width-offset || head.X < offset || head.Y >= height-offset || head.Y < offset {
		return true
	}

	return false
}

func (sb *SnakeBody) GetHead() geometry.Point {
	return sb.Parts[len(sb.Parts)-1]
}
