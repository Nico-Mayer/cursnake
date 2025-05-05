package gamestate

import (
	"slices"
	"strconv"
	"time"
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/fruit"
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
	snakeBody := snake.NewSnakeBody(5, 10, settings.Get().SnakeBodyOptions.StartLength)
	fruitsCollection := fruit.NewFruitCollection(settings.Get().NumberOfFruits, snakeBody, screen)

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
		invalidPoints := slices.Concat(s.snakeBody.Parts, s.fruitsCollection.GetFruitPositions())
		// invalidPoints := append(
		// 	append([]geometry.Point{}, s.snakeBody.Parts...),
		// 	s.fruitsCollection.FruitPositions()...,
		// )
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
	char := unicode.ToLower(event.Rune())
	key := event.Key()

	if key == tcell.KeyUp || char == 'w' || char == 'k' {
		s.snakeBody.Up()
	} else if key == tcell.KeyDown || char == 's' || char == 'j' {
		s.snakeBody.Down()
	} else if key == tcell.KeyLeft || char == 'a' || char == 'h' {
		s.snakeBody.Left()
	} else if key == tcell.KeyRight || char == 'd' || char == 'l' {
		s.snakeBody.Right()
	}
}

func drawCheckerboard(screen tcell.Screen) {
	if !settings.Get().CheckerboardBackground {
		return
	}
	width, height := screen.Size()

	for x := range width {
		for y := range height {
			if (x+y)%2 == 0 {
				screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.Background(tcell.NewRGBColor(20, 20, 20)))
			}
		}
	}
}

func drawGameBorder(screen tcell.Screen) {
	width, height := screen.Size()
	borderStyle := tcell.StyleDefault.Background(tcell.ColorDefault)

	// Draw top border
	for x := range width {
		screen.SetContent(x, 0, '─', nil, borderStyle)
	}
	// Draw bottom border
	for x := range width {
		screen.SetContent(x, height-1, '─', nil, borderStyle)
	}
	// Draw left border
	for y := range height {
		screen.SetContent(0, y, '│', nil, borderStyle)
	}
	// Draw right border
	for y := range height {
		screen.SetContent(width-1, y, '│', nil, borderStyle)
	}

	// Draw corners
	screen.SetContent(0, 0, '┌', nil, borderStyle)              // Top-left
	screen.SetContent(width-1, 0, '┐', nil, borderStyle)        // Top-right
	screen.SetContent(0, height-1, '└', nil, borderStyle)       // Bottom-left
	screen.SetContent(width-1, height-1, '┘', nil, borderStyle) // Bottom-right
}
