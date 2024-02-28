package golang

import (
	"fmt"
	"path/filepath"

	"github.com/aiyengar2/magic/pkg/path"
	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/magefile/mage/mg"
)

var (
	CoverageDir      = "coverage"
	CoverageOutFile  = filepath.Join(CoverageDir, "coverage.out")
	CoverageHTMLFile = filepath.Join(CoverageDir, "coverage.html")
)

func Test() error {
	mg.Deps(setup)

	// Note: for running tests on specific modules, it is encouraged to use Cover instead.
	out, err := cmd.Run(Go, "test", "-cover", "-tags=test", "./...")
	fmt.Printf("%s\n", out)
	return err
}

func Cover(module string) error {
	mg.Deps(setup)

	if len(module) == 0 {
		module = "./..."
	}

	mg.Deps(mg.F(path.Mkdir, CoverageDir))

	err := cmd.RunV(Go, "test", "-v", module, "-covermode=count", fmt.Sprintf("-coverpkg=%s", module), fmt.Sprintf("-coverprofile=%s", CoverageOutFile))
	if err != nil {
		return err
	}
	err = cmd.RunV(Go, "tool", "cover", "-html", CoverageOutFile, "-o", CoverageHTMLFile)
	if err != nil {
		return err
	}
	fmt.Printf("SUCCESS: Run 'open %s' to see the results or refresh the page.\n", CoverageHTMLFile)
	return nil
}
