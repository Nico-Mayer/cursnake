package gamestate

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/fruit"
	"github.com/nico-mayer/cursnake/internal/geometry"
	"github.com/nico-mayer/cursnake/internal/utils"
	"github.com/nico-mayer/cursnake/settings"
	"github.com/nico-mayer/cursnake/snake"
)

type InGameState struct {
	score            int
	snakeBody        *snake.SnakeBody
	fruitsCollection *fruit.FruitCollection
}

func NewInGameState(screen tcell.Screen) *InGameState {
	snakeBody := snake.NewSnakeBody(5, 10, 73)
	fruitsCollection := fruit.NewFruitCollection(settings.GetSettings().NumberOfFruits, snakeBody, screen)

	return &InGameState{
		score:            0,
		snakeBody:        snakeBody,
		fruitsCollection: fruitsCollection,
	}
}

func (s *InGameState) Update(delta time.Duration, screen tcell.Screen) (GameState, bool) {
	width, height := screen.Size()

	var fruitEaten bool
	for _, f := range s.fruitsCollection.Fruits {
		fruitEaten = s.snakeBody.CheckFruitCollision(f.Position)
		if fruitEaten {
			s.score += 10
			var invalidPoints []geometry.Point

			for _, part := range s.snakeBody.Parts {
				invalidPoints = append(invalidPoints, part)
			}

			for _, fruit := range s.fruitsCollection.Fruits {
				invalidPoints = append(invalidPoints, fruit.Position)
			}

			f.Respawn(screen, invalidPoints)
			break
		}
	}

	gameOver := s.snakeBody.CheckSelfCollision()
	if gameOver {
		newState := NewGameOverState()
		return newState, true
	}

	s.snakeBody.Update(delta, width, height, fruitEaten)
	return nil, false
}

func (s *InGameState) Draw(screen tcell.Screen) {
	screen.Clear()
	drawCheckerboard(screen)
	utils.DrawText(1, 1, "Score: "+strconv.Itoa(s.score), screen, tcell.ColorWhite)
	s.snakeBody.Render(screen)
	s.fruitsCollection.Render(screen)
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
				screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.Background(tcell.NewRGBColor(20, 20, 20)))
			}
		}
	}
}
