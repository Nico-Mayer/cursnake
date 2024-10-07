package main

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/fruit"
	"github.com/nico-mayer/cursnake/snake"
)

const frameRate = time.Second / 25

type Game struct {
	Screen    tcell.Screen
	Score     int
	snakeBody *snake.SnakeBody
	fruit     *fruit.Fruit
}

func (g *Game) Run() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	g.Screen.SetStyle(defStyle)
	width, height := g.Screen.Size()
	g.fruit = fruit.NewFruit(width, height, g.snakeBody)
	lastUpdate := time.Now()

	for {
		delta := time.Since(lastUpdate)
		g.Screen.Clear()

		// Update
		fruitEaten := checkFruitCollision(g.snakeBody, g.fruit)
		if fruitEaten {
			g.Score += 10
			g.fruit = fruit.NewFruit(width, height, g.snakeBody)
		}
		g.snakeBody.Update(delta, width, height, fruitEaten)

		// Render
		g.renderScore("Score: " + strconv.Itoa(g.Score))
		g.snakeBody.Render(g.Screen)
		g.fruit.Render(g.Screen)
		g.Screen.Show()

		lastUpdate = time.Now()
		time.Sleep(frameRate)
	}
}

func checkFruitCollision(sb *snake.SnakeBody, fruit *fruit.Fruit) bool {
	head := sb.Parts[len(sb.Parts)-1]

	return head == fruit.Position
}

func (g Game) renderScore(text string) {
	row := 1
	col := 1
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	for _, r := range text {
		g.Screen.SetContent(col, row, r, nil, style)
		col++
	}
}
