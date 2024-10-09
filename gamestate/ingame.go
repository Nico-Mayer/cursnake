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
	"github.com/nico-mayer/cursnake/sound"
)

type InGameState struct {
	score            int
	snakeBody        *snake.SnakeBody
	fruitsCollection *fruit.FruitCollection
}

func NewInGameState(screen tcell.Screen) *InGameState {
	snakeBody := snake.NewSnakeBody(5, 10, 170)
	fruitsCollection := fruit.NewFruitCollection(settings.GetSettings().NumberOfFruits, snakeBody, screen)

	return &InGameState{
		score:            0,
		snakeBody:        snakeBody,
		fruitsCollection: fruitsCollection,
	}
}

func (s *InGameState) Update(delta time.Duration, screen tcell.Screen) GameState {
	width, height := screen.Size()

	fruitEaten, fruit := s.fruitsCollection.CheckCollision(s.snakeBody.GetHead())
	if fruitEaten {
		s.score += 10
		sb := sound.GetManager()
		go sb.Play("pick.mp3")
		invalidPoints := append(
			append([]geometry.Point{}, s.snakeBody.Parts...),
			s.fruitsCollection.FruitPositions()...,
		)
		fruit.Respawn(screen, invalidPoints)
	}

	gameOver := s.snakeBody.CheckSelfCollision() || s.snakeBody.CheckWallCollision(screen)
	if gameOver {
		newState := NewGameOverState()
		return newState
	}

	s.snakeBody.Update(delta, width, height, fruitEaten)
	return nil
}

func (s *InGameState) Draw(screen tcell.Screen) {
	screen.Clear()
	drawCheckerboard(screen)
	utils.DrawText(1, 1, "Score: "+strconv.Itoa(s.score), screen, tcell.ColorWhite)
	s.snakeBody.Render(screen)
	s.fruitsCollection.Render(screen)
	drawGameBorder(screen)
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

func drawGameBorder(screen tcell.Screen) {
	if settings.GetSettings().OpenWalls {
		return
	}
	width, height := screen.Size()
	borderStyle := tcell.StyleDefault.Background(tcell.ColorGray)

	// Draw top border
	for x := 0; x < width; x++ {
		screen.SetContent(x, 0, ' ', nil, borderStyle)
	}
	// Draw bottom border
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-1, ' ', nil, borderStyle)
	}
	// Draw left border
	for y := 0; y < height; y++ {
		screen.SetContent(0, y, ' ', nil, borderStyle)
	}
	// Draw right border
	for y := 0; y < height; y++ {
		screen.SetContent(width-1, y, ' ', nil, borderStyle)
	}
}
