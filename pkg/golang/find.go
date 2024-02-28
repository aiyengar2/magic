package golang

import (
	"fmt"
	"path/filepath"

	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/magefile/mage/mg"
	"golang.org/x/tools/go/packages"
)

func FindModule(path string) (string, error) {
	modules, err := FindAll(path)
	if err != nil {
		return "", err
	}
	switch len(modules) {
	case 0:
		return "", fmt.Errorf("no package exists at %s", path)
	case 1:
		return modules[0], nil
	default:
		return "", fmt.Errorf("found %d local packages, expected 1", len(modules))
	}
}

func FindAll(path string) ([]string, error) {
	mg.Deps(setup)

	// Taken from https://github.com/ko-build/ko/blob/4100cbad34f9480cb24d6a03e4fa1be31b6a19d0/pkg/build/gobuild.go#L204
	dir := filepath.Clean(path)
	if dir == "." {
		dir = ""
	}
	cfg := &packages.Config{
		Mode: packages.NeedModule,
		Dir:  dir,
	}

	cmd.PrintGo("finding all packages in path %s", path)
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, err
	}

	var modules []string
	for _, pkg := range pkgs {
		m := pkg.Module
		if m == nil {
			continue
		}
		modules = append(modules, m.Path)
	}
	return modules, nil
}
