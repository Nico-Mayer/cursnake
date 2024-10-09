package sound

import (
	"bytes"
	"embed"
	"log"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"github.com/nico-mayer/cursnake/settings"
)

type SoundManager struct {
	context *oto.Context
}

//go:embed sounds/*.mp3
var soundFiles embed.FS
var instance *SoundManager
var once sync.Once

func init() {
	once.Do(func() {
		instance = newSoundManager()
	})
}

func GetManager() *SoundManager {
	return instance
}

func newSoundManager() *SoundManager {
	options := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}

	context, readyChan, err := oto.NewContext(options)
	if err != nil {
		log.Fatalf("Failed to create oto context: %v", err)
	}
	<-readyChan

	return &SoundManager{context: context}
}

func (sm *SoundManager) Play(fileName string) {
	if settings.GetSettings().Mute {
		return
	}

	data, err := soundFiles.ReadFile("sounds/" + fileName)
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}
	reader := bytes.NewReader(data)
	decodedMp3, err := mp3.NewDecoder(reader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}
	player := sm.context.NewPlayer(decodedMp3)
	defer player.Close()

	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Microsecond)
	}
}
