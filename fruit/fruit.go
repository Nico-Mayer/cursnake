package fruit

import (
	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/internal/geometry"
	"github.com/nico-mayer/cursnake/internal/utils"
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
	x, y := generateNewPos(screen, invalidPoints)

	return &Fruit{
		Position: geometry.Point{X: x, Y: y},
		Style:    tcell.StyleDefault.Foreground(tcell.ColorRed),
	}
}

func (fc *FruitCollection) Render(screen tcell.Screen) {
	for _, fruit := range fc.Fruits {
		screen.SetContent(fruit.Position.X, fruit.Position.Y, '\u25CF', nil, fruit.Style)
	}
}

func (fc *FruitCollection) FruitPositions() []geometry.Point {
	var positions []geometry.Point
	for _, fruit := range fc.Fruits {
		positions = append(positions, fruit.Position)
	}
	return positions
}

func (fc *FruitCollection) CheckCollision(snakeHead geometry.Point) (bool, *Fruit) {
	for _, fruit := range fc.Fruits {
		if fruit.Position == snakeHead {
			return true, fruit
		}
	}
	return false, nil
}

func (f *Fruit) Respawn(screen tcell.Screen, invalidPoints []geometry.Point) {
	x, y := generateNewPos(screen, invalidPoints)

	f.Position.X = x
	f.Position.Y = y
}

func generateNewPos(screen tcell.Screen, invalidPoints []geometry.Point) (x, y int) {
	width, height := screen.Size()
	x = utils.RandRange(1, width-1)
	y = utils.RandRange(1, height-1)

	for checkIfInvalid(x, y, invalidPoints) {
		x = utils.RandRange(1, width-1)
		y = utils.RandRange(1, height-1)
	}
	return x, y
}

func checkIfInvalid(x, y int, invalidPoints []geometry.Point) bool {
	for _, point := range invalidPoints {
		if point.X == x && point.Y == y {
			return true
		}
	}
	return false
}
