package utils

import (
	"math/rand"
	"strconv"
	"strings"

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

func Clamp(min, value, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func RandRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func HexColorToInt32(hexColor string) (int32, error) {
	hexColor = strings.TrimPrefix(hexColor, "#")

	parsedValue, err := strconv.ParseInt(hexColor, 16, 32)
	if err != nil {
		return 0, err
	}

	return int32(parsedValue), nil
}
