//go:build debug_mac

package audio

import (
	"mamela/buildconstraints"
	"mamela/merror"
	"os"
	"path/filepath"
)

func init() {

	exPath, err := filepath.Abs(filepath.Dir(os.Args[0])) //get the current working directory
	if err != nil {
		merror.ShowError("", err)
	}

	dir := exPath

	LibDir = dir + buildconstraints.PathSeparator + "lib" + buildconstraints.PathSeparator + "mac"
	LibExt = ".dylib"
	NotifyInitReady <- true
}
