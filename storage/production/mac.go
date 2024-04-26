//go:build prod_mac

package storage

import (
	"mamela/buildConstraints"
	"mamela/merror"
	"mamela/storage"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		merror.ShowError("", e)
	}
	exPath := filepath.Dir(ex)
	exPathArr := strings.Split(exPath, buildConstraints.PathSeparator)
	exPathArr = exPathArr[0 : len(exPathArr)-1]

	dir := strings.Join(exPathArr, buildConstraints.PathSeparator) + buildConstraints.PathSeparator + "db"

	storage.StorageFile = dir + buildConstraints.PathSeparator + "data.json"
}
