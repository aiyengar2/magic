package magic

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aiyengar2/magic/pkg/docker"
	"github.com/aiyengar2/magic/pkg/git"
	"github.com/aiyengar2/magic/pkg/golang"
	"github.com/aiyengar2/magic/pkg/golangci_lint"
	"github.com/aiyengar2/magic/pkg/tool"
	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/aiyengar2/magic/pkg/utils/env"
	magicversion "github.com/aiyengar2/magic/pkg/version"
	"github.com/magefile/mage/mg"
)

const (
	CrossEnvVar = "CROSS"
)

var (
	Go = golang.BuildOptions{
		Targets: []string{"."},
		Platforms: []string{
			"linux/amd64",
			"linux/arm64",
			"darwin/amd64",
			"darwin/arm64",
			"windows/amd64",
		},
		OtherLDFlags: env.OSDependent[[]string]{
			Windows: []string{},
			Linux:   []string{"-extldflags", "-static", "-s"},
		}.Resolve(),
	}
	Docker = docker.BuildOptions{{
		Dockerfile: filepath.Join("package", "Dockerfile"),
		Image:      docker.NewImage(golang.AppName(".")),
	}}
)

// version

var (
	VersionFn = versionFn
)

// Version (params: none) is the version detected by this build system
func Version() {
	mg.SerialDeps(VersionFn)
	cmd.PrintFinish("version")
}

func versionFn() error {
	fmt.Printf("binary version: %s\n", magicversion.GetVersion())
	fmt.Printf("package version: %s\n", docker.NewImage(golang.AppName(".")))
	return nil
}

// build

var (
	BuildFn = buildFn
)

// Build (params: none) builds all your Go modules for your machine
func Build() {
	mg.SerialDeps(BuildFn)
	cmd.PrintFinish("build")
}

func buildFn() error {
	if len(os.Getenv(CrossEnvVar)) > 0 {
		Go.Platforms = []string{}
	}
	return golang.BuildCustom(Go)
}

// cover

var (
	CoverFn = coverFn
)

// Cover (params: <module:string>) produces a coverage.html detailing go coverage for each of your modules
func Cover(path string) {
	mg.SerialDeps(mg.F(CoverFn, path))
	cmd.PrintFinish("cover")
}

func coverFn(path string) error {
	mg.SerialDeps(Build, mg.F(golang.Cover, path))
	return nil
}

// test

var (
	TestFn = testFn
)

// Test (params: none) runs go tests on your modules
func Test() {
	mg.SerialDeps(TestFn)
	cmd.PrintFinish("test")
}

func testFn() error {
	mg.SerialDeps(Build, golang.Test)
	return nil
}

// validate

var (
	ValidateFn = validateFn
)

// Validate (params: none) runs validation steps on your repository
func Validate() {
	mg.SerialDeps(ValidateFn)
	cmd.PrintFinish("validate")
}

func validateFn() error {
	mg.SerialDeps(
		Build,
		golang.Generate,
		golang.Fmt,
	)
	err := tool.Exists(golangci_lint.GolangCILint)
	if err == nil {
		mg.Deps(golangci_lint.Run)
	} else {
		cmd.PrintWarning("%s not found, skipping", golangci_lint.GolangCILint)
	}

	dirty := git.Dirty()
	if dirty {
		_ = git.Status()
		_ = git.Diff()
		return fmt.Errorf("git is dirty")
	}
	return nil
}

// package

var (
	PackageFn = packageFn
)

// Package (params: none) packages your Go programs into Docker images on your machine
func Package() {
	mg.SerialDeps(PackageFn)
	cmd.PrintFinish("package")
}

func packageFn() error {
	mg.SerialDeps(Build)
	return docker.Build(Docker, len(os.Getenv(CrossEnvVar)) > 0)
}

// publish

var (
	PublishFn = publishFn
)

// Publish (params: none) publishes your Docker images to a configured repository
func Publish() {
	mg.SerialDeps(PublishFn)
	cmd.PrintFinish("publish")
}

func publishFn() error {
	mg.SerialDeps(Build)
	return docker.Publish(Docker)
}

// default

var (
	DefaultFn = defaultFn
)

// Default (params: none) runs build, test, package
func Default() {
	mg.SerialDeps(DefaultFn)
	cmd.PrintFinish("default")
}

func defaultFn() error {
	mg.SerialDeps(Build, Test, Package)
	return nil
}

// ci

var (
	CIFn = ciFn
)

// CI (params: none) runs buildCross, test, validate, packageCross
func CI() {
	mg.SerialDeps(CIFn)
	cmd.PrintFinish("ci")
}

func ciFn() error {
	mg.SerialDeps(Build, Test, Validate, Package)
	return nil
}

// release

var (
	ReleaseFn = releaseFn
)

// Release (params: none) runs ci
func Release() {
	mg.SerialDeps(ReleaseFn)
	cmd.PrintFinish("release")
}

func releaseFn() error {
	mg.SerialDeps(Build, Test, Validate, Publish)
	return nil
}
