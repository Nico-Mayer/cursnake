package gamestate

import (
	"time"

	"github.com/gdamore/tcell"
)

type GameState interface {
	Update(time.Duration, tcell.Screen) (newState GameState)
	Draw(tcell.Screen)
	HandleInput(*tcell.EventKey)
}
