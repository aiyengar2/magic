package golang

import (
	"fmt"
	"path/filepath"
)

const (
	UnknownAppName = "UNKNOWN"
)

func AppName(target string) string {
	gomod, err := FindModule(target)
	if err != nil {
		return UnknownAppName
	}
	return filepath.Base(gomod)
}

func BinaryFmt(app, goos, goarch string) string {
	var ext string
	if goos == "windows" {
		ext = ".exe"
	}
	return fmt.Sprintf("%s-%s-%s%s", app, goos, goarch, ext)
}
