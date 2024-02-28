package path

import (
	"io/fs"
	"os"

	"github.com/aiyengar2/magic/pkg/utils/cmd"
)

func Mkdir(path string) error {
	cmd.PrintGo("creating directory at path '%s'", path)
	return os.MkdirAll(path, fs.FileMode(0750))
}
