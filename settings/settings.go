package settings

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed default-settings.json
var defaultSettingsData []byte

type CursnakeSettings struct {
	TargetFPS              *int  `json:"TargetFPS"`
	Sound                  *bool `json:"Sound"`
	CheckerboardBackground *bool `json:"CheckerboardBackground"`
	NumberOfFruits         *int  `json:"NumberOfFruits"`
	OpenWalls              *bool `json:"OpenWalls"`
}

var settings *CursnakeSettings

func init() {
	settings = newCursnakeSettings()
}

func GetSettings() *CursnakeSettings {
	return settings
}

func newCursnakeSettings() *CursnakeSettings {
	defaultSetting := loadDefaultSettings()
	userSettings := loadUserSettings()

	mergedSettings := mergeSettings(defaultSetting, userSettings)

	// Windows-specific defaults
	if runtime.GOOS == "windows" {
		mergedSettings.CheckerboardBackground = boolPointer(false)
		mergedSettings.TargetFPS = intPointer(25)
	}

	return mergedSettings
}

func loadDefaultSettings() *CursnakeSettings {
	var defaultConfig CursnakeSettings
	if err := json.Unmarshal(defaultSettingsData, &defaultConfig); err != nil {
		panic("failed to load default settings")
	}

	return &defaultConfig
}

func loadUserSettings() *CursnakeSettings {
	userSettingsPath := filepath.Join(os.Getenv("HOME"), ".config", "cursnake", "settings.json")
	userSettingsFile, err := os.ReadFile(userSettingsPath)
	if err != nil {
		return &CursnakeSettings{}
	}

	var userSettings CursnakeSettings
	if err := json.Unmarshal(userSettingsFile, &userSettings); err != nil {
		return &CursnakeSettings{}
	}

	return &userSettings
}

func mergeSettings(defaultSettings, userSettings *CursnakeSettings) *CursnakeSettings {
	merged := &CursnakeSettings{}

	if userSettings.TargetFPS != nil {
		merged.TargetFPS = userSettings.TargetFPS
	} else {
		merged.TargetFPS = defaultSettings.TargetFPS
	}

	if userSettings.NumberOfFruits != nil {
		merged.NumberOfFruits = userSettings.NumberOfFruits
	} else {
		merged.NumberOfFruits = defaultSettings.NumberOfFruits
	}

	if userSettings.Sound != nil {
		merged.Sound = userSettings.Sound
	} else {
		merged.Sound = defaultSettings.Sound
	}

	if userSettings.CheckerboardBackground != nil {
		merged.CheckerboardBackground = userSettings.CheckerboardBackground
	} else {
		merged.CheckerboardBackground = defaultSettings.CheckerboardBackground
	}

	if userSettings.OpenWalls != nil {
		merged.OpenWalls = userSettings.OpenWalls
	} else {
		merged.OpenWalls = defaultSettings.OpenWalls
	}

	return merged
}

func intPointer(i int) *int {
	return &i
}

func boolPointer(b bool) *bool {
	return &b
}
