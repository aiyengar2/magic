package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aiyengar2/magic/pkg/path"
	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/aiyengar2/magic/pkg/utils/env"
	"github.com/magefile/mage/mg"
)

type BuildOptions struct {
	Targets        []string
	Platforms      []string
	BuildOverrides map[string]string
	OtherLDFlags   []string
}

// BuildCustom a go module with custom 'go build' arguments.
func BuildCustom(opts BuildOptions) error {
	if len(opts.Targets) == 0 {
		cmd.PrintWarning("no go modules to build.")
		return nil
	}

	mg.Deps(setup)

	buildInfos := make([]struct {
		Path    string
		App     string
		LDFlags []string
	}, len(opts.Targets))

	for i, target := range opts.Targets {
		// set path
		buildInfos[i].Path = target

		// compute gomod
		gomod, err := FindModule(target)
		if err != nil {
			return err
		}

		// set app
		buildInfos[i].App = filepath.Base(gomod)

		// set ldflags
		var ldFlags []string
		for k, v := range opts.BuildOverrides {
			ldFlags = append(ldFlags, fmt.Sprintf("-X %s/pkg/%s=%s", gomod, k, v))
		}
		buildInfos[i].LDFlags = append(ldFlags, opts.OtherLDFlags...)
	}

	// create the bin directory
	mg.Deps(mg.F(path.Mkdir, "bin"))

	// run the builds
	for _, buildInfo := range buildInfos {
		for _, osArchVersion := range env.GetOSArchVersions(opts.Platforms...) {
			mg.Deps(mg.F(build, buildInfo.App, "bin", buildInfo.Path, strings.Join(buildInfo.LDFlags, " "), osArchVersion.OS, osArchVersion.Arch))
		}
	}

	return nil
}

func build(app string, binDir string, path string, ldFlagsStr, os string, arch string) error {
	goos := os
	goarch := arch
	env := map[string]string{
		"GOOS":        goos,
		"GOARCH":      goarch,
		"CGO_ENABLED": "0",
	}
	args := []string{"build"}
	args = append(args, "-ldflags", ldFlagsStr)
	args = append(args, "-o", filepath.Join(binDir, BinaryFmt(app, goos, goarch)))
	args = append(args, path)
	return cmd.RunWithV(env, Go, args...)
}
