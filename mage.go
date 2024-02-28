//go:build mage

package main

import (
	"github.com/aiyengar2/magic/pkg/docker"

	// mage:import
	"github.com/aiyengar2/magic/pkg/magic"
)

func init() {
	// # Go Settings

	// Use this to override the platforms you want to build
	// magic.Go.Platforms = []string{
	// 	"linux/amd64",
	// 	"linux/arm64",
	// 	"darwin/amd64",
	// 	"darwin/arm64",
	// 	"windows/amd64",
	// }

	// Use this to override the targets for where to build the go program
	// magic.Go.Targets = []string{
	// 	".",
	// }

	// Use this to override variables on build time
	// magic.Go.BuildOverrides = map[string]string{
	// 	"version.Version":   version.Version,
	// 	"version.GitCommit": version.Commit,
	// }

	// # Docker Settings

	// Option 1: If you have multiple Dockerfiles, override magic.Docker entirely
	// magic.Docker = docker.BuildOptions{
	// 	{
	// 		Dockerfile: filepath.Join("package", "Dockerfile"),
	// 		Tag:        version.GetTag("magic"),
	// 		PlatformTargets: []docker.PlatformTarget{
	// 			{
	// 				PlatformVersions: []string{"linux/amd64", "linux/arm64"},
	// 				BuildArgs: map[string]string{
	// 					"BASE_IMAGE": "registry.suse.com/bci/bci-micro:15.4.17.1",
	// 				},
	// 			},
	// 			{
	// 				PlatformVersions: []string{"windows/amd64/ltsc2019"},
	// 				BuildArgs: map[string]string{
	// 					"BASE_IMAGE": "mcr.microsoft.com/windows/nanoserver:ltsc2019",
	// 				},
	// 			},
	// 			{
	// 				PlatformVersions: []string{"windows/amd64/ltsc2022"},
	// 				BuildArgs: map[string]string{
	// 					"BASE_IMAGE": "mcr.microsoft.com/windows/nanoserver:ltsc2022",
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// Option 2: If you have a single Dockerfile, just override the default
	// magic.Docker[0].Dockerfile = filepath.Join("package", "Dockerfile")
	// magic.Docker[0].Tag = version.GetTag(golang.AppName("."))

	// Here is an example of how to set up your platform targets based on the Dockerfile
	// in package/Dockerfile. By default, you can build this with `mage package` and push
	// a multi-arch image to a repository with `REPO=<my-fork> mage publish`
	magic.Docker[0].PlatformTargets = []docker.PlatformTarget{
		{
			TargetBuildStage: "linux",
			PlatformVersions: []string{"linux/amd64"},
			BuildArgs: map[string]string{
				"BASE_IMAGE": "registry.suse.com/bci/bci-micro:15.4.17.1",
				"OS":         "linux",
				"ARCH":       "amd64",
			},
		},
		{
			TargetBuildStage: "linux",
			PlatformVersions: []string{"linux/arm64"},
			BuildArgs: map[string]string{
				"BASE_IMAGE": "registry.suse.com/bci/bci-micro:15.4.17.1",
				"OS":         "linux",
				"ARCH":       "arm64",
			},
		},
		{
			TargetBuildStage: "windows",
			PlatformVersions: []string{"windows/amd64/ltsc2019"},
			BuildArgs: map[string]string{
				"BASE_IMAGE": "mcr.microsoft.com/windows/nanoserver:ltsc2019",
				"OS":         "windows",
				"ARCH":       "amd64",
			},
		},
		{
			TargetBuildStage: "windows",
			PlatformVersions: []string{"windows/amd64/ltsc2022"},
			BuildArgs: map[string]string{
				"BASE_IMAGE": "mcr.microsoft.com/windows/nanoserver:ltsc2022",
				"OS":         "windows",
				"ARCH":       "amd64",
			},
		},
	}

	// # Overriding Default Build Targets

	// If you would like to override any of the default targets
	// magic.BuildFn = func() error {
	// 	// custom logic
	// 	return nil
	// }

	// Or you can add onto it
	// buildFn := magic.BuildFn
	// magic.BuildFn = func() error {
	// 	buildFn()
	// 	// custom logic here
	// 	return nil
	// }
}
