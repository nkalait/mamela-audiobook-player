package ui

import (
	"bytes"
	"image"
	"mamela/audio"
	"mamela/merror"
	"mamela/types"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var bookArt *canvas.Image
var bookTitle *widget.Label
var chapterTitle *widget.Label
var bookFullLength *canvas.Text
var playingPosition *canvas.Text
var playerButtonPlay *widget.Button
var playerButtonPause *widget.Button

var playerButtonFastRewind *widget.Button
var playerButtonFastForward *widget.Button
var playerButtonSkipNext *widget.Button
var playerButtonSkipPrevious *widget.Button

func createPlayingLayout() *fyne.Container {
	initUI()
	hideUIItems()

	containerTop := container.NewVBox(bookTitle, bookFullLength)
	containerPositionDetails := container.NewVBox(
		playingPosition,
		layoutPlayerButtons(),
	)
	playingVBox := container.NewPadded(container.NewBorder(
		containerTop,
		containerPositionDetails,
		nil,
		nil,
		container.NewBorder(nil, chapterTitle, nil, nil, bookArt),
	))

	a := container.NewBorder(
		nil,
		nil,
		nil,
		createVolumeSlider(),
		playingVBox,
	)

	allBorder := container.NewBorder(nil, playTimeScrubber, nil, nil, a)

	go func() {
		for playingBook := range audio.UpdateNowPlayingChannel {
			if bookTitle.Hidden {
				showUIItems()
			}
			// If playingBook is empty
			if playingBook.Path == "" {
				clearCurrentlyPlaying()
			} else {
				updatePlaying(playingBook)
			}
		}
	}()
	return allBorder
}

func initUI() {
	initTitle()
	initChapterTitle()
	initBookArt()
	initPlayingPosition()
	initPlayerButtons()
	initFullBookLength()
	initPlayTimeScrubberSlider()
}

func showUIItems() {
	bookArt.Show()
	bookTitle.Show()
	bookFullLength.Show()
	playingPosition.Show()

	chapterTitle.Show()
	playerButtonSkipPrevious.Show()
	playerButtonFastRewind.Show()
	playerButtonPause.Show()
	playerButtonPlay.Show()
	playerButtonFastForward.Show()
	playerButtonSkipNext.Show()

	showVolumeSlider()
	showPlayTimeScrubber()
}

func hideUIItems() {
	bookArt.Hide()
	bookTitle.Hide()
	bookFullLength.Hide()
	playingPosition.Hide()

	chapterTitle.Hide()
	playerButtonSkipPrevious.Hide()
	playerButtonFastRewind.Hide()
	playerButtonPause.Hide()
	playerButtonPlay.Hide()
	playerButtonFastForward.Hide()
	playerButtonSkipNext.Hide()

	hideVolumeSlider()
	hidePlayTimeScrubber()
}

func initBookArt() {
	bookArt = canvas.NewImageFromImage(nil)
	bookArt.FillMode = canvas.ImageFillContain
}

func initTitle() {
	bookTitle = widget.NewLabel("")
	bookTitle.Wrapping = fyne.TextWrapBreak
	bookTitle.TextStyle.Bold = true
	bookTitle.Alignment = fyne.TextAlignCenter
}

func initChapterTitle() {
	chapterTitle = widget.NewLabel("")
	chapterTitle.Wrapping = fyne.TextWrapWord
	chapterTitle.Alignment = fyne.TextAlignCenter
}

func initFullBookLength() {
	bookFullLength = canvas.NewText("", theme.ForegroundColor())
	bookFullLength.TextSize = 18
	bookFullLength.Alignment = fyne.TextAlignCenter
}

func initPlayingPosition() {
	playingPosition = canvas.NewText("", theme.ForegroundColor())
	playingPosition.TextSize = 18
	playingPosition.Alignment = fyne.TextAlignCenter
}

func initPlayerButtons() {
	playerButtonPlay = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		audio.Play()
	})
	playerButtonPause = widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {
		audio.Pause()
	})
	playerButtonFastRewind = widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), func() {
		audio.FastRewind()
	})
	playerButtonFastForward = widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		audio.FastForward()
	})
	playerButtonSkipNext = widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
		audio.SkipNext()
	})
	playerButtonSkipPrevious = widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		audio.SkipPrevious()
	})
}

func layoutPlayerButtons() *fyne.Container {
	layout := container.NewHBox(
		playerButtonSkipPrevious,
		playerButtonFastRewind,
		playerButtonPause,
		playerButtonPlay,
		playerButtonFastForward,
		playerButtonSkipNext,
	)
	vbox := container.NewVBox(layout)
	return container.NewCenter(vbox)
}
func clearCurrentlyPlaying() {
	clearBookArt()
	updateTitle("")
	updateChapterTitle("")
	clearPlayingPosition()
	updateBookFullLength("")
	hideUIItems()
}

func updatePlaying(p types.PlayingBook) {
	updateTitle(p.Title)
	updateChapterTitle(p.Chapters[p.CurrentChapter].FileName)
	updatePlayingPosition()
	d := time.Duration(math.Round(p.FullLengthSeconds * 1000000000))
	updateBookFullLength(audio.SecondsToTimeText(d))
	if audio.GetState() == audio.PAUSED || audio.GetState() == audio.STOPPED {
		playerButtonPlay.Show()
		playerButtonPause.Hide()
	} else if audio.GetState() == audio.PLAYING {
		playerButtonPlay.Hide()
		playerButtonPause.Show()
	}
}

func clearBookArt() {
	bookArt.Image = nil
	bookArt.Refresh()
}

func updateBookArt(picBytes []byte) {
	clearBookArt()
	img, _, err := image.Decode(bytes.NewReader(picBytes))
	if err != nil {
		merror.ShowError("Problem loading audio book image", err)
		return
	}
	bookArt.Image = img
	bookArt.Refresh()
}

func updateTitle(title string) {
	bookTitle.Text = cases.Title(language.English).String(title)
	bookTitle.Refresh()
}

func updateChapterTitle(title string) {
	chapterTitle.Text = cases.Title(language.English).String(title)
	chapterTitle.Refresh()
}

func updateBookFullLength(bookLength string) {
	bookFullLength.Text = bookLength
	bookFullLength.Refresh()
}

func updatePlayingPosition() {
	playingPosition.Text = audio.SecondsToTimeText(audio.GetCurrentBookPlayTime())
	playingPosition.Refresh()
}

func clearPlayingPosition() {
	playingPosition.Text = ""
	playingPosition.Refresh()
}
