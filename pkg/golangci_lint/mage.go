package golangci_lint

import (
	"fmt"

	"github.com/aiyengar2/magic/pkg/tool"
	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/aiyengar2/magic/pkg/utils/env"
	"github.com/magefile/mage/mg"
)

var (
	GolangCILint = env.CrossPlatformBinary("golangci-lint")
)

func setup() {
	mg.Deps(mg.F(tool.Exists, GolangCILint))
}

func Run() error {
	mg.Deps(setup)

	out, err := cmd.Run(GolangCILint, "run")
	if err != nil {
		// an exception for golanglint-ci, which pipes errors to stdout
		fmt.Printf("%s\n", out)
	}
	return nil
}
