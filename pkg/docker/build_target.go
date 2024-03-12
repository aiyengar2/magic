package docker

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/aiyengar2/magic/pkg/utils/env"
	"github.com/aiyengar2/magic/pkg/version"
)

var (
	Org = "arvindiyengar" // TODO: modify this when finalizing the repository
	Tag = "v0.0.0-dev"
)

func init() {
	// Get org
	if org := os.Getenv("ORG"); len(org) > 0 {
		Org = org
	}

	// Get tag
	Tag = os.Getenv("TAG")
	if len(Tag) == 0 {
		Tag = version.Version
	}
}

type BuildTarget struct {
	Dockerfile      string
	Image           ImageName
	PlatformTargets []PlatformTarget

	Context   string
	BuildArgs map[string]string
	OtherArgs []string
}

func NewImage(repo string) ImageName {
	return ImageName{
		Org:  Org,
		Repo: repo,
		Tag:  Tag,
	}
}

type ImageName struct {
	Org  string
	Repo string
	Tag  string
}

func (n ImageName) String() string {
	return fmt.Sprintf("%s/%s:%s", n.Org, n.Repo, n.Tag)
}

type PlatformTarget struct {
	TargetBuildStage string
	BuildArgs        map[string]string

	PlatformVersions []string
}

func (t BuildTarget) build(publish bool) error {
	if len(t.PlatformTargets) == 0 {
		cmd.PrintWarning("no platforms to build for %s.", t.Image)
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
			platformTagSlice := []string{t.Image.String()}
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
	return pushManifest(t.Image.String(), builtTags...)
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
		cmd.PrintWarning("inferring PlatformTarget to match %s for %s, this should be defined in the BuildOptions", osArchs, t.Image)
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
				fmt.Printf("Skip: building platform %s for image %s when CROSS is not set\n", p, t.Image)
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
