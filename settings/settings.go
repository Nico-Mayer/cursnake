package settings

import (
	"runtime"
	"time"
)

type CursnakeSettings struct {
	Os                     string
	FrameRate              time.Duration
	Sound                  bool
	CheckerboardBackground bool
	MultipleFruits         bool
	OpenWalls              bool
}

var settings *CursnakeSettings

func init() {
	settings = newCursnakeSettings()
}

func newCursnakeSettings() *CursnakeSettings {
	system := runtime.GOOS
	frameRate := time.Second / 60
	sound := true
	checkerboardBG := true
	multipleFruits := false
	openWalls := true

	if system == "windows" {
		checkerboardBG = false
		frameRate = time.Second / 25
	}

	return &CursnakeSettings{
		Os:                     system,
		FrameRate:              frameRate,
		Sound:                  sound,
		CheckerboardBackground: checkerboardBG,
		MultipleFruits:         multipleFruits,
		OpenWalls:              openWalls,
	}
}

func GetSettings() *CursnakeSettings {
	return settings
}
