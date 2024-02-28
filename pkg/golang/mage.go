package golang

import (
	"fmt"

	"github.com/aiyengar2/magic/pkg/tool"
	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/aiyengar2/magic/pkg/utils/env"
	"github.com/magefile/mage/mg"
)

var (
	Go = env.CrossPlatformBinary("go")
)

func setup() {
	mg.Deps(mg.F(tool.Exists, Go))
}

func Generate() error {
	mg.Deps(setup)

	return cmd.RunV(Go, "generate")
}

func Fmt() error {
	mg.Deps(setup)

	out, err := cmd.Run(Go, "fmt", "./...")
	if len(out) > 0 {
		fmt.Printf("%s\n", out)
		if err == nil {
			return fmt.Errorf("gofmt failed")
		}
	}
	return err
}
