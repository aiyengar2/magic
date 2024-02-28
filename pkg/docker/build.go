package docker

import (
	"fmt"
	"strings"

	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/magefile/mage/mg"
)

type BuildOptions []BuildTarget

func Build(buildTargets BuildOptions, cross bool) error {
	if len(buildTargets) == 0 {
		cmd.PrintWarning("no docker images to build.")
		return nil
	}

	mg.Deps(setup)

	return buildOrPublish(false, buildTargets, cross)
}

func Publish(buildTargets BuildOptions) error {
	if len(buildTargets) == 0 {
		cmd.PrintWarning("no docker images to publish.")
		return nil
	}

	mg.Deps(setup)

	return buildOrPublish(true, buildTargets, true)
}

func buildOrPublish(publish bool, buildTargets BuildOptions, cross bool) error {
	if err := BuildxCreate(); err != nil {
		return err
	}
	defer func() {
		err := BuildxRemove()
		if err != nil {
			panic(err)
		}
	}()

	for _, buildTarget := range buildTargets {
		err := buildTarget.filterEnvironments(cross).build(publish)
		if err != nil {
			return err
		}
	}
	return nil
}

func build(dockerfile string, dockerContext string, tag string, targetBuildStage string, platform string, buildArgs string, otherArgs string, publish bool) error {
	if !publish {
		// see if we can build this
		supported := false
		for _, supportedPlatform := range DefaultPlatforms {
			if platform == supportedPlatform {
				supported = true
			}
		}
		if !supported {
			fmt.Printf("Skip: loading a built image of platform %s is not supported on this machine\n", platform)
			return nil
		}
	}

	args := []string{"buildx", "build", "--pull"}
	args = append(args, "--file", dockerfile)
	args = append(args, "--tag", tag)

	if publish {
		args = append(args, "--push")
	} else {
		args = append(args, "--load")
	}
	if len(targetBuildStage) > 0 {
		args = append(args, "--target", targetBuildStage)
	}
	if len(platform) > 0 {
		args = append(args, "--platform", platform)
	}
	if len(buildArgs) > 0 {
		var buildArgsSlice []string
		for _, buildArg := range strings.Split(buildArgs, "\n") {
			buildArgsSlice = append(buildArgsSlice, "--build-arg")
			buildArgsSlice = append(buildArgsSlice, buildArg)
		}
		args = append(args, buildArgsSlice...)
	}
	if len(otherArgs) > 0 {
		args = append(args, strings.Split(otherArgs, " ")...)
	}
	if len(dockerContext) == 0 {
		dockerContext = "."
	}
	args = append(args, dockerContext)
	return cmd.RunV(Docker, args...)
}

func pushManifest(manifestTag string, tags ...string) error {
	args := []string{"buildx", "imagetools", "create"}
	args = append(args, "--tag", manifestTag)
	args = append(args, tags...)
	return cmd.RunV(Docker, args...)
}
