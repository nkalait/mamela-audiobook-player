package ui

import (
	"fmt"
	"mamela/storage"
	"testing"
)

func TestLeftPaneHeaderText(t *testing.T) {
	BuildUI("appLabel", true)
	fmt.Println(storage.Data.BookList)
	header := "Load Books     "
	if bookListHeaderTxt.Text != header {
		t.Errorf("bookListHeaderTxt should be \"%s\"", header)
	}
}
