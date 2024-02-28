package docker

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/aiyengar2/magic/pkg/utils/env"
)

type BuildTarget struct {
	Dockerfile      string
	Tag             string
	PlatformTargets []PlatformTarget

	Context   string
	BuildArgs map[string]string
	OtherArgs []string
}

type PlatformTarget struct {
	TargetBuildStage string
	BuildArgs        map[string]string

	PlatformVersions []string
}

func (t BuildTarget) build(publish bool) error {
	if len(t.PlatformTargets) == 0 {
		cmd.PrintWarning("no platforms to build for %s.", t.Tag)
		return nil
	}

	var builtTags []string
	for _, pt := range t.PlatformTargets {
		// build args
		var buildArgs []string
		buildArgs = append(buildArgs, cmd.MapToEnvSlice(t.BuildArgs)...)
		buildArgs = append(buildArgs, cmd.MapToEnvSlice(pt.BuildArgs)...)
		sort.Strings(buildArgs)
		buildArgsStr := strings.Join(buildArgs, "\n")

		// other args
		sort.Strings(t.OtherArgs)
		otherArgsStr := strings.Join(t.OtherArgs, " ")

		// run for each platform
		sort.Strings(pt.PlatformVersions)
		for _, platformVersion := range pt.PlatformVersions {
			oav := env.GetOSArchVersions(platformVersion)[0]
			platformTagSlice := []string{t.Tag}
			for _, attr := range []string{oav.OS, oav.Arch, oav.Version} {
				if len(attr) > 0 {
					platformTagSlice = append(platformTagSlice, attr)
				}
			}
			platformTag := strings.Join(platformTagSlice, "-")
			err := build(t.Dockerfile, t.Context, platformTag, pt.TargetBuildStage, oav.OSArch(), buildArgsStr, otherArgsStr, publish)
			if err != nil {
				return err
			}
			builtTags = append(builtTags, platformTag)
		}
	}
	if !publish {
		return nil
	}
	return pushManifest(t.Tag, builtTags...)
}

func (t BuildTarget) filterEnvironments(cross bool) BuildTarget {
	if cross {
		return t
	}

	// If not cross, assume that you only build on the first build platform (default)
	osArchs := DefaultPlatforms[:1]

	if len(t.PlatformTargets) == 0 {
		t.PlatformTargets = []PlatformTarget{{
			PlatformVersions: osArchs,
		}}
		cmd.PrintWarning("inferring PlatformTarget to match %s for %s, this should be defined in the BuildOptions", osArchs, t.Tag)
	}

	// get the list of osArchs
	osArchMap := make(map[string]bool)
	for _, osArch := range osArchs {
		osArchMap[osArch] = true
	}

	var platformTargets []PlatformTarget
	for _, pt := range t.PlatformTargets {
		var platforms []string
		for _, p := range pt.PlatformVersions {
			if osArchMap[p] {
				platforms = append(platforms, p)
			} else {
				fmt.Printf("Skip: building platform %s for image %s when CROSS is not set\n", p, t.Tag)
			}
		}
		if len(platforms) > 0 {
			pt.PlatformVersions = platforms
			platformTargets = append(platformTargets, pt)
		}
	}
	t.PlatformTargets = platformTargets
	return t
}
