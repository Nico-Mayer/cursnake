package fruit

import (
	"math/rand"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/internal/geometry"
	"github.com/nico-mayer/cursnake/snake"
)

type FruitCollection struct {
	Fruits []*Fruit
}

type Fruit struct {
	Position geometry.Point
	Style    tcell.Style
}

func NewFruitCollection(size int, sb *snake.SnakeBody, screen tcell.Screen) *FruitCollection {
	var fruits []*Fruit

	var invalidPoints []geometry.Point

	for _, part := range sb.Parts {
		invalidPoints = append(invalidPoints, part)
	}

	for i := 0; i < size; i++ {
		fruits = append(fruits, newFruit(screen, invalidPoints))
		invalidPoints = append(invalidPoints, fruits[i].Position)
	}

	return &FruitCollection{
		Fruits: fruits,
	}
}

func newFruit(screen tcell.Screen, invalidPoints []geometry.Point) *Fruit {
	width, height := screen.Size()
	xPos := rand.Intn(width)
	yPos := rand.Intn(height)

	for checkIfInvalid(xPos, yPos, invalidPoints) {
		xPos = rand.Intn(width)
		yPos = rand.Intn(height)
	}

	return &Fruit{
		Position: geometry.Point{X: xPos, Y: yPos},
		Style:    tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed),
	}
}

func (fc *FruitCollection) Render(screen tcell.Screen) {
	for _, fruit := range fc.Fruits {
		screen.SetContent(fruit.Position.X, fruit.Position.Y, '\u25CF', nil, fruit.Style)
	}
}

func (f *Fruit) Respawn(screen tcell.Screen, invalidPoints []geometry.Point) {
	width, height := screen.Size()
	xPos := rand.Intn(width)
	yPos := rand.Intn(height)

	for checkIfInvalid(xPos, yPos, invalidPoints) {
		xPos = rand.Intn(width)
		yPos = rand.Intn(height)
	}

	f.Position.X = xPos
	f.Position.Y = yPos
}

func checkIfInvalid(x, y int, invalidPoints []geometry.Point) bool {
	for _, point := range invalidPoints {
		if point.X == x && point.Y == y {
			return true
		}
	}

	return false
}
