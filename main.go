package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell"
	"github.com/nico-mayer/cursnake/snake"
)

func main() {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	snakeBody := snake.NewSnakeBody(5, 10, 4)

	game := Game{
		Screen:    screen,
		snakeBody: snakeBody,
	}

	go game.Run()

	for {
		switch event := game.Screen.PollEvent().(type) {
		case *tcell.EventResize:
			game.Screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				game.Screen.Fini()
				os.Exit(0)
			} else if event.Key() == tcell.KeyUp || event.Rune() == 'w' {
				game.snakeBody.Up()
			} else if event.Key() == tcell.KeyDown || event.Rune() == 's' {
				game.snakeBody.Down()
			} else if event.Key() == tcell.KeyLeft || event.Rune() == 'a' {
				game.snakeBody.Left()
			} else if event.Key() == tcell.KeyRight || event.Rune() == 'd' {
				game.snakeBody.Right()
			}
		}
	}
}
