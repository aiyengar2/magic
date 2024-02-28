package env

import (
	"fmt"
	"path/filepath"
)

var (
	Newline = OSDependent[string]{
		Linux:   "\n",
		Windows: "\r\n",
	}.Resolve()
)

type OSDependent[T interface{}] struct {
	Windows T
	Linux   T
}

func (d OSDependent[T]) Resolve() T {
	if OS == "windows" {
		return d.Windows
	}
	return d.Linux
}

func CrossPlatformBinary(binary string) string {
	if len(filepath.Ext(binary)) > 0 {
		panic(fmt.Errorf("cannot use env.CrossPlatformBinary on binary with extension: %s", binary))
	}
	return OSDependent[string]{
		Linux:   binary,
		Windows: fmt.Sprintf("%s.exe", binary),
	}.Resolve()
}
