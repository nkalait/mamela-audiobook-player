package ui

import (
	"mamela/audio"
	"mamela/bundled"
	"mamela/storage"
	"mamela/ui/customtheme"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/sqweek/dialog"
)

var Testing = false

var (
	mamelaApp  fyne.App
	MainWindow fyne.Window
)

func BuildUI(appLabel string, testing bool) {
	Testing = testing
	mamelaApp = app.New()
	if !testing {
		setupTheming()
		setSystemTrayMenu()
	}
	mainContainer := arrangeUI()
	prepareMainWindow(appLabel, mainContainer)
}

func setSystemTrayMenu() {
	if desk, ok := mamelaApp.(desktop.App); ok {
		m := fyne.NewMenu("Mamela",
			fyne.NewMenuItem("Pause", func() {
				audio.Pause()
			}),
			fyne.NewMenuItem("Play", func() {
				audio.Play()
			}),
			fyne.NewMenuItem("FastForward", func() {
				audio.FastForward()
			}),
			fyne.NewMenuItem("Rewind", func() {
				audio.FastRewind()
			}),
			fyne.NewMenuItem("Next Chapter", func() {
				audio.SkipNext()
			}),
			fyne.NewMenuItem("Previous Chapter", func() {
				audio.SkipPrevious()
			}),
		)
		var res = &fyne.StaticResource{
			StaticName:    "",
			StaticContent: bundled.ResourceIconPng.StaticContent,
		}
		desk.SetSystemTrayIcon(res)
		desk.SetSystemTrayMenu(m)
	}
}

func prepareMainWindow(label string, c *fyne.Container) {
	MainWindow = mamelaApp.NewWindow(label)
	MainWindow.SetContent(c)
	if !Testing {
		MainWindow.Resize(fyne.NewSize(800, 600))
		MainWindow.SetMainMenu(makeMainMenu())
		MainWindow.SetCloseIntercept(func() {
			audio.ExitListener <- true
			time.Sleep(1 * time.Second)
			MainWindow.Close()
			time.Sleep(1 * time.Second)
			os.Exit(0) // this is to also quit the system tray menu
		})
		MainWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
			if k.Name == "Space" {
				if audio.GetState() == audio.PLAYING {
					audio.Pause()
				} else {
					audio.Play()
				}
			} else if k.Name == "Up" || k.Name == "Down" {
				adjustVolumeOnKeyPress(string(k.Name))
			} else if k.Name == "Left" || k.Name == "Right" {
				adjustPlayTimeScrubberOnKeyPress(string(k.Name))
			} else if k.Name == "S" || k.Name == "s" {
				audio.Stop()
			}
		})
		MainWindow.CenterOnScreen()
		MainWindow.ShowAndRun()
	}
}

func makeMainMenu() *fyne.MainMenu {
	menuItem := fyne.NewMenuItem("Root Folder", func() {
		openSelectRootFolderDialog()
	})
	menu := fyne.NewMenu("File", menuItem)
	mainMenu := fyne.NewMainMenu(menu)
	return mainMenu
}

func openSelectRootFolderDialog() {
	path, err := dialog.Directory().Title("Open root folder").Browse()
	if err != nil {
		dialog.Message(err.Error())
	} else if path != "" {
		notifyRootFolderSelected <- true
		refreshButton.Hide()
		storage.ClearBooks()
		audio.ClearCurrentlyPlaying()
		storage.Data.Root = path
		refreshBookList()
		storage.SaveDataToStorageFile()
	}
}

func setupTheming() {
	if customtheme.IsLight() {
		mamelaApp.Settings().SetTheme(customtheme.LightTheme())
	} else {
		mamelaApp.Settings().SetTheme(customtheme.DarkTheme())
	}
}

func arrangeUI() *fyne.Container {
	// UI to show list of audio books
	leftPane := container.NewStack(
		canvas.NewRectangle(customtheme.GetColour(customtheme.ColourNameBackgroundLight)),
		container.NewPadded(container.NewBorder(
			generateBookListContainerTop(),
			nil, nil, nil,
			initBookList(),
		)),
	)

	// Part of UI to place UI elements pertaining to currently playing audio book
	currentlyPlayingContainer := createPlayingLayout()

	// Place all parts UI parts together
	bodyParts := container.NewBorder(
		nil,
		nil,
		leftPane,
		nil,
		currentlyPlayingContainer,
	)

	// bodyParts := container.NewHSplit(leftPane, currentlyPlayingContainer)
	// bodyParts.Offset = 0
	return container.NewGridWithColumns(1, bodyParts)
}
