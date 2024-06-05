package ui

import (
	"mamela/audio"
	"mamela/storage"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var volumeSlider *widget.Slider

const maxVolume = float64(10000)

var currentVolume = float64(0)

func init() {
	volumeSlider = widget.NewSlider(currentVolume, maxVolume) // 10000 is bass's Stream global volume max level
}

func createVolumeSlider() *fyne.Container {
	volumeSlider.Orientation = 1 // vertical
	onDragVolumeSlider(volumeSlider)
	hBox := container.NewHBox(volumeSlider, container.NewPadded())
	return container.NewPadded(hBox)
}

func onDragVolumeSlider(slider *widget.Slider) {
	slider.OnChangeEnded = func(f float64) {
		currentVolume = f
		audio.NotifyVolumeSliderDragged <- currentVolume
		storage.SaveVolumeLevel(currentVolume)
	}
}

func hideVolumeSlider() {
	volumeSlider.Hide()
}

func showVolumeSlider() {
	volumeSlider.Show()
	if !Testing {
		volumeSlider.SetValue(storage.GetVolumeLevel())
		audio.NotifyVolumeSliderDragged <- storage.GetVolumeLevel()
	}
}

const volumeAdjustValue = float64(500)

func adjustVolumeOnKeyPress(keyName string) {
	switch keyName {
	case "Up":
		if volumeAdjustValue+currentVolume > maxVolume {
			currentVolume = maxVolume
		} else {
			currentVolume += volumeAdjustValue
		}
	case "Down":
		if currentVolume-volumeAdjustValue < 0 {
			currentVolume = 0
		} else {
			currentVolume -= volumeAdjustValue
		}
	}
	audio.NotifyVolumeSliderDragged <- currentVolume
	volumeSlider.SetValue(currentVolume)
}
