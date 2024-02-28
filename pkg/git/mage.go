package git

import (
	"fmt"

	"github.com/aiyengar2/magic/pkg/tool"
	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/aiyengar2/magic/pkg/utils/env"
	"github.com/magefile/mage/mg"
)

var (
	Git = env.CrossPlatformBinary("git")
)

func setup() {
	mg.Deps(mg.F(tool.Exists, Git))
}

func Status() error {
	mg.Deps(setup)

	out, err := cmd.Run(Git, "status")
	fmt.Printf("%s\n", out)
	return err
}

func Diff() error {
	mg.Deps(setup)

	out, err := cmd.Run(Git, "diff")
	fmt.Printf("%s\n", out)
	return err
}
