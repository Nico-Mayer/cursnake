package gamestate

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/fruit"
	"github.com/nico-mayer/cursnake/internal/utils"
	"github.com/nico-mayer/cursnake/settings"
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

func (s *InGameState) Update(delta time.Duration, screen tcell.Screen) (GameState, bool) {
	width, height := screen.Size()

	fruitEaten := s.snakeBody.CheckFruitCollision(s.fruit.Position)
	gameOver := s.snakeBody.CheckSelfCollision()
	if gameOver {
		newState := NewGameOverState()
		return newState, true
	}

	if fruitEaten {
		s.score += 10
		s.fruit = fruit.NewFruit(width, height, s.snakeBody)
	}
	s.snakeBody.Update(delta, width, height, fruitEaten)
	return nil, false
}

func (s *InGameState) Draw(screen tcell.Screen) {
	screen.Clear()
	drawCheckerboard(screen)
	utils.DrawText(1, 1, "Score: "+strconv.Itoa(s.score), screen, tcell.ColorWhite)
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

func drawCheckerboard(screen tcell.Screen) {
	if !settings.GetSettings().CheckerboardBackground {
		return
	}
	width, height := screen.Size()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if (x+y)%2 == 0 {
				screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.Background(tcell.NewRGBColor(10, 10, 10)))
			}
		}
	}
}
