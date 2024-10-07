package main

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/gamestate"
)

type Cursnake struct {
	screen           tcell.Screen
	currentGameState gamestate.GameState
}

const (
	frameRate = time.Second / 144
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer screen.Fini()

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	cursnake := &Cursnake{
		screen:           screen,
		currentGameState: gamestate.NewInGameState(screen),
	}

	go cursnake.Run()

	for {
		switch event := cursnake.screen.PollEvent().(type) {
		case *tcell.EventResize:
			cursnake.screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				screen.Fini()
				os.Exit(0)
			}
			cursnake.currentGameState.HandleInput(event)
		}
	}
}

func (g *Cursnake) Run() {
	lastUpdate := time.Now()

	for {
		g.currentGameState.Update(time.Since(lastUpdate), g.screen)
		lastUpdate = time.Now()

		g.currentGameState.Draw(g.screen)

		time.Sleep(frameRate)
	}
}
