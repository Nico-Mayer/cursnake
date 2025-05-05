package main

import (
	"log"
	"os"
	"time"
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/gamestate"
	"github.com/nico-mayer/cursnake/settings"
)

type Cursnake struct {
	screen           tcell.Screen
	currentGameState gamestate.GameState
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer screen.Fini()

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
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC || unicode.ToLower(event.Rune()) == 'q' {
				screen.Fini()
				os.Exit(0)
			}
			cursnake.currentGameState.HandleInput(event)
		}
	}
}

func (cursnake *Cursnake) Run() {
	frameTime := time.Second / time.Duration(settings.Get().TargetFPS)
	lastUpdate := time.Now()

	// Gameloop
	for {
		delta := time.Since(lastUpdate)
		newState := cursnake.currentGameState.Update(delta, cursnake.screen)
		if newState != nil {
			cursnake.currentGameState = newState
		}
		lastUpdate = time.Now()

		// Draw
		cursnake.currentGameState.Draw(cursnake.screen)

		time.Sleep(frameTime)
	}
}
