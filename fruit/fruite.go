package fruit

import (
	"math/rand"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/geometry"
	"github.com/nico-mayer/cursnake/snake"
)

type Fruit struct {
	Position geometry.Point
	style    tcell.Style
}

func NewFruit(width, height int, sb *snake.SnakeBody) *Fruit {
	xPos := rand.Intn(width)
	yPos := rand.Intn(height)

	for checkIfInsideSnake(xPos, yPos, sb) {
		xPos = rand.Intn(width)
		yPos = rand.Intn(height)
	}

	return &Fruit{
		Position: geometry.Point{X: xPos, Y: yPos},
		style:    tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed),
	}
}

func (f Fruit) Render(screen tcell.Screen) {
	screen.SetContent(f.Position.X, f.Position.Y, '\u25CF', nil, f.style)
}

func checkIfInsideSnake(x, y int, sb *snake.SnakeBody) bool {
	for _, part := range sb.Parts {
		if part.X == x && part.Y == y {
			return true
		}
	}

	return false
}
