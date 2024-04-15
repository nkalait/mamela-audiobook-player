package audio

import (
	"fmt"
	"mamela/buildconstraints"
	"mamela/err"
	"mamela/types"
	"time"

	bass "github.com/pteich/gobass"
)

// Event listeners
var (
	exitListener = make(chan bool) // for stopping to listen to playing events
	exitAudio    = make(chan bool) // for unloading audio stuff
)

const (
	Stopped = iota
	Paused
	Playing
)

var ChannelAudioState = make(chan int)
var UpdateNowPlayingChannel = make(chan types.PlayingBook)

// Holds data structures important to playing an audiobook
var player Player

// Initiate Bass
func init() {
	go func() {
		plugins := loadPlugins()
		defer func(p []uint32) {
			tearDown(p)
		}(plugins)
		initBass()
		<-exitAudio
	}()
}

// Unload loadded Bass plugins and free all resources used by Bass
func tearDown(plugins []uint32) {
	for _, p := range plugins {
		bass.PluginFree(p)
	}
	bass.Free()
}

// Initialise Bass
func initBass() {
	e := bass.Init(-1, 44100, bass.DeviceStereo, 0, 0)
	err.PanicError(e)
	bass.SetVolume(100)
}

// Load pluggins needed by Bass
func loadPlugins() []uint32 {
	pluginLibbassAac, e := bass.PluginLoad("lib"+buildconstraints.PathSeparator+"mac"+buildconstraints.PathSeparator+"libbass_aac.dylib", bass.StreamDecode)
	err.PanicError(e)
	pluginLibbassOpus, e := bass.PluginLoad("lib"+buildconstraints.PathSeparator+"mac"+buildconstraints.PathSeparator+"libbassopus.dylib", bass.StreamDecode)
	err.PanicError(e)

	plugins := make([]uint32, 2)
	plugins = append(plugins, pluginLibbassAac)
	plugins = append(plugins, pluginLibbassOpus)

	return plugins
}

const TickerDuration = 500 * time.Millisecond

var Ticker = time.NewTicker(TickerDuration)

// Start listening to audio playing event and exit event
func StartChannelListener(exitApp chan bool) {
	player.updater = UpdateNowPlayingChannel
	go func() {
	RoutineLoop:
		for {
			select {
			// case <-time.After(time.Second):
			case <-Ticker.C:
				updateUICurrentlyPlayingInfo()
			case <-exitListener:
				break RoutineLoop
			}
		}
		exitAudio <- true
		time.Sleep(time.Second * 2)
		exitApp <- true
	}()
}

// Pad number below 10 with a zero
func pad(i int) string {
	s := fmt.Sprint(i)
	if i < 10 {
		s = "0" + fmt.Sprint(i)
	}
	return s
}

// Convert seconds to time in hh : mm : ss
func SecondsToTimeText(seconds time.Duration) string {
	var h int = int(seconds.Seconds()) / 3600
	var m int = (int(seconds.Seconds()) / 60) % 60
	var s int = int(seconds.Seconds()) % 60

	sh := pad(h)
	sm := pad(m)
	ss := pad(s)

	return sh + " : " + sm + " : " + ss
}

func GetCurrentBookPlayingDuration(p types.PlayingBook) time.Duration {
	pos := p.Position
	if p.CurrentChapter > 0 {
		for i := p.CurrentChapter - 1; i >= 0; i-- {
			pos = pos + time.Duration(p.Chapters[i].LengthInSeconds*1000000000)
		}
	}
	return pos
}

// Update the currently playing audio book information on the UI
func updateUICurrentlyPlayingInfo() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		err.PanicError(e)
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_STOPPED {
			bytePosition, e := player.channel.GetPosition(bass.POS_BYTE)
			err.PanicError(e)
			p, e := player.channel.Bytes2Seconds(bytePosition)
			err.PanicError(e)

			currentlyAt := player.currentBook.Position.Round(time.Second)
			skipAt := time.Duration(player.currentBook.Chapters[player.currentBook.CurrentChapter].LengthInSeconds * 1000000000).Round(time.Second)
			if currentlyAt == skipAt {
				skipToNextFile(&player, true)
			}

			posInWholeBook := GetCurrentBookPlayingDuration(player.currentBook).Round(time.Second)
			wholeBookLength := time.Duration(player.currentBook.FullLengthSeconds * 1000000000).Round(time.Second)
			if posInWholeBook == wholeBookLength {
				player.currentBook.Finished = true
				Ticker.Stop()
				ChannelAudioState <- Stopped
			}
			var d time.Duration = time.Duration(p * 1000000000)
			player.currentBook.Position = time.Duration(d)
		} else if active == bass.ACTIVE_STOPPED {
			if !player.currentBook.Finished {
				player.currentBook.Position = 0
			}
			ChannelAudioState <- Stopped
		}
		player.updater <- player.currentBook
	}
}

type UpdateFolderArtCallBack func(playingBook types.PlayingBook)

func LoadAndPlay(playingBook types.PlayingBook, updaterFolderArtCallback UpdateFolderArtCallBack) {
	// c, e := bass.StreamCreateURL("http://music.myradio.ua:8000/PopRock_news128.mp3", bass.DeviceStereo)
	player.currentBook = playingBook
	stopPlayingIfPlaying(player.channel, player)

	chapter := player.currentBook.CurrentChapter
	e := loadAudioBookFile(player.currentBook.FullPath + buildconstraints.PathSeparator + player.currentBook.Chapters[chapter].FileName)
	if e == nil {
		startPlaying()
	}

	if updaterFolderArtCallback != nil {
		updaterFolderArtCallback(player.currentBook)
	}
	updateUICurrentlyPlayingInfo()
}

func stopPlayingIfPlaying(c bass.Channel, p Player) {
	if c != 0 {
		a, e := c.IsActive()
		err.PanicError(e)
		if a == bass.ACTIVE_PLAYING || a == bass.ACTIVE_PAUSED {
			p.stop()
		}
	}
}

func loadAudioBookFile(fullPath string) error {
	var e error = nil
	player.channel, e = bass.StreamCreateFile(fullPath, 0, bass.AsyncFile)
	if e != nil {
		err.ShowError("There seems to be a problem loading the the audio book file(s)", e)
	}

	return e
}

func startPlaying() error {
	e := player.channel.SetPosition(0, bass.POS_BYTE)
	if e != nil {
		err.ShowError("There seems to be a problem playing the the audio book file(s)", e)
	} else {
		player.play()
	}
	return e
}

func skipToNextFile(p *Player, forceSkip bool) bool {
	skipped := false
	if p.channel != 0 {
		active, e := p.channel.IsActive()
		err.ShowError("Error skipping to next chapter", e)
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_PAUSED || forceSkip {
			numChapters := len(p.currentBook.Chapters)
			if numChapters > 0 {
				if p.currentBook.CurrentChapter < numChapters-1 {
					p.currentBook.CurrentChapter = p.currentBook.CurrentChapter + 1
					LoadAndPlay(p.currentBook, nil)
					skipped = true
				}
			}
		}
	}
	return skipped
}

func skipToPreviousFile(p *Player) bool {
	skipped := false
	if p.channel != 0 {
		active, e := p.channel.IsActive()
		err.ShowError("Error skipping to previous chapter", e)
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_PAUSED {
			numChapters := len(p.currentBook.Chapters)
			if numChapters > 0 {
				if p.currentBook.CurrentChapter > 0 {
					p.currentBook.CurrentChapter = p.currentBook.CurrentChapter - 1
					LoadAndPlay(p.currentBook, nil)
					skipped = true
				} else {
					e = p.channel.SetPosition(0, bass.POS_BYTE)
					if e != nil {
						err.ShowError("Error to skipping to start", e)
					}
				}
			}
		}
	}
	return skipped
}
