package utils

import (
	"github.com/gdamore/tcell"
)

func DrawText(row, col int, text string, screen tcell.Screen, color tcell.Color) {
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(color)
	for _, r := range text {
		screen.SetContent(col, row, r, nil, style)
		col++
	}
}

func DrawTextCenter(row, col int, text string, screen tcell.Screen, color tcell.Color) {
	col = col - len(text)/2
	DrawText(row, col, text, screen, color)
}
