package gamestate

import (
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/internal/utils"
)

type GameOverState struct {
	cursorPos      int
	selectedOption string
	options        []string
}

func NewGameOverState() *GameOverState {
	return &GameOverState{
		cursorPos: 0,
		options:   []string{"Restart", "Exit"},
	}
}

func (s *GameOverState) Update(delta time.Duration, screen tcell.Screen) (GameState, bool) {
	if len(s.selectedOption) != 0 {
		switch s.selectedOption {
		case "Restart":
			{
				newState := NewInGameState(screen)
				return newState, true
			}
		case "Exit":
			{
				screen.Fini()
				os.Exit(0)
			}
		}
	}
	return nil, false
}

func (s *GameOverState) Draw(screen tcell.Screen) {
	width, height := screen.Size()
	screen.Clear()
	utils.DrawTextCenter(height/2, width/2, "Game Over", screen, tcell.ColorWhite)

	for i, option := range s.options {
		color := tcell.ColorWhite
		y := height/2 + 2 + i*2 // 2 rows of vertical spacing between options
		x := width / 2

		if i == s.cursorPos {
			option = "> " + option + " <"
			color = tcell.ColorGreenYellow
		}

		utils.DrawTextCenter(y, x, option, screen, color)
	}

	screen.Show()
}

func (s *GameOverState) HandleInput(event *tcell.EventKey) {
	if event.Key() == tcell.KeyUp || event.Rune() == 'w' {
		if s.cursorPos > 0 {
			s.cursorPos--
		}
	} else if event.Key() == tcell.KeyDown || event.Rune() == 's' {
		if s.cursorPos <= 0 {
			s.cursorPos++
		}
	} else if event.Key() == tcell.KeyEnter || event.Rune() == ' ' {
		s.selectedOption = s.options[s.cursorPos]
	}
}
