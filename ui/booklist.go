package ui

import (
	"fmt"
	"image/color"
	"mamela/audio"
	"mamela/buildconstraints"
	"mamela/helpers"
	"mamela/storage"
	"mamela/types"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
)

var notifyRootFolderSelected = make(chan bool)
var loadingBooksTicker = time.NewTicker(500 * time.Millisecond)

var refreshButton *widget.Button
var bookListVBox *fyne.Container
var bookListHeaderTxt *canvas.Text

func init() {
	loadingBooksTicker.Stop()
}

// Initialise part of the UI that lists audio books and listen to update events
func initBookList() *fyne.Container {
	bookListVBox = container.New(layout.NewVBoxLayout())
	bookListContainer := initBookPane(bookListVBox)
	// Note, passing in true here will overwrite the book store in the storage file on disk,
	// this also means losing saved book last play positions
	updateBookList(false)

	// Listen to book list update events
	go func() {
		for update := range audio.UpdateBookListChannel {
			if update {
				audio.Stop()
				audio.ClearPlayer()
				clearCurrentlyPlaying()
				storage.ClearBooks()
				audio.NotifyBookPlayTimeSliderDragged <- -1
				updateBookList(true)
				notifyRootFolderSelected <- false
			}
		}
	}()

	dots := 1
	go func() {
		for range loadingBooksTicker.C {
			switch dots {
			case 1:
				bookListHeaderTxt.Text = "Loading."
			case 2:
				bookListHeaderTxt.Text = "Loading.."
			case 3:
				bookListHeaderTxt.Text = "Loading..."
			}
			bookListHeaderTxt.Refresh()
			dots++
			if dots == 4 {
				dots = 1
			}
		}
	}()
	go func() {
		for runAnimation := range notifyRootFolderSelected {
			if runAnimation {
				loadingBooksTicker.Reset(500 * time.Millisecond)
			} else {
				loadingBooksTicker.Stop()
				updateBookListHeader()
				refreshButton.Show()
			}
		}
	}()

	return bookListContainer
}

func updateBookListHeader() {
	if len(storage.Data.BookList) > 0 {
		bookListHeaderTxt.Text = "Loaded Books"
	} else {
		// TODO find a better we of doing the padding below
		bookListHeaderTxt.Text = "Load Books     "
	}
	bookListHeaderTxt.Refresh()
}

func generateBookListContainerTop() *fyne.Container {
	bookListHeaderTxt = canvas.NewText("", theme.ForegroundColor())
	bookListHeaderTxt.TextSize = 24
	bookListHeaderTxt.TextStyle.Bold = true
	spacer := canvas.NewText("    ", color.Transparent)
	refreshButton = createRefreshButton()
	top := container.NewHBox(bookListHeaderTxt, spacer, container.NewVBox(refreshButton))
	return top
}

func initBookPane(bookListVBox *fyne.Container) *fyne.Container {
	// TODO the dots below are just to give the scroller the desired width, NEED TO FIND A WAY TO DO THIS BETTER!!
	dots := canvas.NewText("..........................................................", color.Transparent)
	bookListScroller := container.NewVScroll(bookListVBox)
	hBox := container.NewStack(dots, bookListScroller)
	return hBox
}

func createRefreshButton() *widget.Button {
	icon := theme.ViewRefreshIcon()
	button := widget.NewButtonWithIcon("", icon, func() {
		notifyRootFolderSelected <- true
		refreshButton.Hide()
		refreshBookList()
	})
	return button
}

func refreshBookList() {
	audio.UpdateBookListChannel <- true
}

// Update the part of UI showing list of audio books,
// if readRootFolder is true then we scan the root folder and save contents
// to the storage file
func updateBookList(readRootFolder bool) {
	if storage.Data.Root != "" {
		bookListVBox.RemoveAll()
		bookListVBox.Refresh()
		if readRootFolder {
			parseRootFolder()
		}
		for _, book := range storage.Data.BookList {
			if book.Path == "" {
				fmt.Println("!!!!!!!!!!!1=================")
			}
			// if book.Path != "" {
			book.Metadata = getFileMetaData(book)
			bookTileLayout := NewMyListItemWidget(book)
			bookListVBox.Add(bookTileLayout)
			loadPreviousBookOnLoad(book.Path, bookTileLayout.Button)
			// }
		}
		bookListVBox.Refresh()
	}
	updateBookListHeader()
}

func loadPreviousBookOnLoad(bookFolderName string, bookLoadButton *widget.Button) {
	if storage.Data.CurrentBookFolder == bookFolderName {
		go func() {
			for {
				time.Sleep(500 * time.Millisecond)
				if MainWindow != nil {
					bookLoadButton.OnTapped()
					audio.Pause()
					break
				}
			}
		}()
	}
}

func parseRootFolder() {
	if storage.Data.Root == "" {
		return
	}
	rootFolderEntries := helpers.ReadRootDirectory()
	bookList := helpers.GetBookList(rootFolderEntries)
	storage.SaveBookListToStorageFile(bookList)
}

func getBookFile(b types.Book) *os.File {
	path := storage.Data.Root + buildconstraints.PathSeparator + b.Path + buildconstraints.PathSeparator + b.Chapters[0].FileName
	f, _ := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	return f
}

func getFileMetaData(b types.Book) tag.Metadata {
	var meta tag.Metadata = nil
	f := getBookFile(b)
	defer f.Close()
	if f != nil {
		meta, _ = tag.ReadFrom(f)
	}
	return meta
}
