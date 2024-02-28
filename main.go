package main

import (
	"fmt"

	"github.com/aiyengar2/magic/pkg/version"
)

func main() {
	fmt.Printf("binary version: %s\n", version.GetVersion())
	fmt.Printf("package version: %s\n", version.GetTag("magic"))
}
