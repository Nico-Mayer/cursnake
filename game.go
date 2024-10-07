package main

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/fruit"
	"github.com/nico-mayer/cursnake/snake"
)

const (
	targetFPS  = 144
	frameTime  = time.Second / targetFPS
	updateTime = time.Second / 60
)

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
	lastRender := time.Now()

	for {
		now := time.Now()

		elapsedUpdate := now.Sub(lastUpdate)
		elapsedRender := now.Sub(lastRender)

		if elapsedUpdate >= updateTime {
			g.Screen.Clear()
			fruitEaten := checkFruitCollision(g.snakeBody, g.fruit)
			if fruitEaten {
				g.Score += 10
				g.fruit = fruit.NewFruit(width, height, g.snakeBody)
			}
			g.snakeBody.Update(time.Since(lastUpdate), width, height, fruitEaten)
			lastUpdate = now
		}

		if elapsedRender >= frameTime {
			g.renderScore("Score: " + strconv.Itoa(g.Score))
			g.snakeBody.Render(g.Screen)
			g.fruit.Render(g.Screen)
			g.Screen.Show()
			lastRender = now
		}

		time.Sleep(time.Until(lastRender.Add(frameTime)))
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
