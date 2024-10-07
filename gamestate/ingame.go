package gamestate

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/fruit"
	"github.com/nico-mayer/cursnake/snake"
)

type InGameState struct {
	score     int
	snakeBody *snake.SnakeBody
	fruit     *fruit.Fruit
}

func NewInGameState(screen tcell.Screen) *InGameState {
	width, height := screen.Size()
	snakeBody := snake.NewSnakeBody(5, 10, 4)

	return &InGameState{
		score:     0,
		snakeBody: snakeBody,
		fruit:     fruit.NewFruit(width, height, snakeBody),
	}
}

func (s *InGameState) Update(delta time.Duration, screen tcell.Screen) {
	width, height := screen.Size()

	fruitEaten := s.snakeBody.CheckFruitCollision(s.fruit.Position)
	if fruitEaten {
		s.score += 10
		s.fruit = fruit.NewFruit(width, height, s.snakeBody)
	}
	s.snakeBody.Update(delta, width, height, fruitEaten)
}

func (s *InGameState) Draw(screen tcell.Screen) {
	screen.Clear()
	renderScore("Score: "+strconv.Itoa(s.score), screen)
	s.snakeBody.Render(screen)
	s.fruit.Render(screen)
	screen.Show()
}

func (s *InGameState) HandleInput(event *tcell.EventKey) {
	if event.Key() == tcell.KeyUp || event.Rune() == 'w' {
		s.snakeBody.Up()
	} else if event.Key() == tcell.KeyDown || event.Rune() == 's' {
		s.snakeBody.Down()
	} else if event.Key() == tcell.KeyLeft || event.Rune() == 'a' {
		s.snakeBody.Left()
	} else if event.Key() == tcell.KeyRight || event.Rune() == 'd' {
		s.snakeBody.Right()
	}
}

func renderScore(text string, s tcell.Screen) {
	row := 1
	col := 1
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}
