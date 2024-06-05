package ui

import (
	"testing"
)

func TestLeftPaneHeaderText(t *testing.T) {
	BuildUI("appLabel", true)
	header := "Load Books     "
	if bookListHeaderTxt.Text != header {
		t.Errorf("bookListHeaderTxt should be \"%s\"", header)
	}
}

func TestHideUIItems(t *testing.T) {
	BuildUI("appLabel", true)
	showUIItems()
	hideUIItems()
	if bookArt.Hidden != true {
		t.Errorf("bookArt should be hidden")
	}
	if bookTitle.Hidden != true {
		t.Errorf("bookTitle should be hidden")
	}
	if bookFullLength.Hidden != true {
		t.Errorf("bookFullLength should be hidden")
	}
	if playingPosition.Hidden != true {
		t.Errorf("playingPosition should be hidden")
	}
	if chapterTitle.Hidden != true {
		t.Errorf("chapterTitle should be hidden")
	}
	if playerButtonSkipPrevious.Hidden != true {
		t.Errorf("playerButtonSkipPrevious should be hidden")
	}
	if playerButtonFastRewind.Hidden != true {
		t.Errorf("playerButtonFastRewind should be hidden")
	}
	if playerButtonPause.Hidden != true {
		t.Errorf("playerButtonPause should be hidden")
	}
	if playerButtonPlay.Hidden != true {
		t.Errorf("playerButtonPlay should be hidden")
	}
	if playerButtonFastForward.Hidden != true {
		t.Errorf("playerButtonFastForward should be hidden")
	}
	if playerButtonSkipNext.Hidden != true {
		t.Errorf("playerButtonSkipNext should be hidden")
	}
	if volumeSlider.Hidden != true {
		t.Errorf("volumeSlider should be hidden")
	}
	if playTimeScrubber.Hidden != true {
		t.Errorf("playTimeScrubber should be hidden")
	}
}
func TestShowUIItems(t *testing.T) {
	BuildUI("appLabel", true)
	showUIItems()
	if bookArt.Hidden != false {
		t.Errorf("bookArt should not be hidden")
	}
	if bookTitle.Hidden != false {
		t.Errorf("bookArt should not be hidden")
	}
	if bookFullLength.Hidden != false {
		t.Errorf("bookTitle should not be hidden")
	}
	if playingPosition.Hidden != false {
		t.Errorf("playingPosition should not be hidden")
	}
	if chapterTitle.Hidden != false {
		t.Errorf("chapterTitle should not be hidden")
	}
	if playerButtonSkipPrevious.Hidden != false {
		t.Errorf("playerButtonSkipPrevious should not be hidden")
	}
	if playerButtonFastRewind.Hidden != false {
		t.Errorf("playerButtonFastRewind should not be hidden")
	}
	if playerButtonPause.Hidden != false {
		t.Errorf("playerButtonPause should not be hidden")
	}
	if playerButtonPlay.Hidden != false {
		t.Errorf("playerButtonPlay should not be hidden")
	}
	if playerButtonFastForward.Hidden != false {
		t.Errorf("playerButtonFastForward should not be hidden")
	}
	if playerButtonSkipNext.Hidden != false {
		t.Errorf("playerButtonSkipNext should not be hidden")
	}
	if volumeSlider.Hidden != false {
		t.Errorf("volumeSlider should not be hidden")
	}
	if playTimeScrubber.Hidden != false {
		t.Errorf("playTimeScrubber should not be hidden")
	}
}
