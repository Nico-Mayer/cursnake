package settings

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

//go:embed default-settings.json
var defaultSettingsData []byte
var once sync.Once

type CursnakeSettings struct {
	TargetFPS              int              `json:"targetFPS"`
	Mute                   bool             `json:"mute"`
	CheckerboardBackground bool             `json:"checkerboardBackground"`
	NumberOfFruits         int              `json:"numberOfFruits"`
	OpenWalls              bool             `json:"openWalls"`
	SnakeBodyOptions       SnakeBodyOptions `json:"snakeBody"`
}

type SnakeBodyOptions struct {
	StartLength int    `json:"startLength"`
	HeadRune    string `json:"headRune"`
	BodyRune    string `json:"bodyRune"`
	Foreground  string `json:"foreground"`
	Background  string `json:"background"`
}

var settings *CursnakeSettings
var defaultSettings *CursnakeSettings

func init() {
	once.Do(func() {
		settings = newCursnakeSettings()
		defaultSettings = loadDefaultSetting()
	})
}

func GetSettings() *CursnakeSettings {
	return settings
}
func GetDefaultSettings() *CursnakeSettings {
	return defaultSettings
}

func newCursnakeSettings() *CursnakeSettings {
	userSettingsFile := loadUserSettings()

	var data1, data2 map[string]interface{}

	if err := json.Unmarshal(defaultSettingsData, &data1); err != nil {
		log.Fatalf("Error unmarshalling file1: %v", err)
	}
	if err := json.Unmarshal(userSettingsFile, &data2); err != nil {
		data2 = data1
	}

	mergedData := mergeMaps(data1, data2)

	mergedJSON, err := json.MarshalIndent(mergedData, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling merged data: %v", err)
	}

	var mergedSettings CursnakeSettings
	if err := json.Unmarshal(mergedJSON, &mergedSettings); err != nil {
		log.Fatalf("Error unmarshalling into struct: %v", err)
	}

	// Windows-specific defaults
	if runtime.GOOS == "windows" {
		mergedSettings.CheckerboardBackground = false
		mergedSettings.TargetFPS = 25
	}

	return &mergedSettings
}

func loadUserSettings() []byte {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return defaultSettingsData
	}
	userSettingsPath := filepath.Join(homeDir, ".config", "cursnake", "settings.json")
	userSettingsFile, err := os.ReadFile(userSettingsPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
		return defaultSettingsData
	}
	return userSettingsFile
}

func loadDefaultSetting() *CursnakeSettings {
	var settings *CursnakeSettings

	if err := json.Unmarshal(defaultSettingsData, &settings); err != nil {
		log.Fatalf("Error unmarshalling into struct: %v", err)
	}

	// Windows-specific defaults
	if runtime.GOOS == "windows" {
		settings.CheckerboardBackground = false
		settings.TargetFPS = 25
	}

	return settings
}

func mergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	for k, v := range map2 {
		map1[k] = v
	}
	return map1
}
