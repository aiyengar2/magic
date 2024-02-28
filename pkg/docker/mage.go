package docker

import (
	"fmt"
	"strings"

	"github.com/aiyengar2/magic/pkg/tool"
	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/aiyengar2/magic/pkg/utils/env"
	"github.com/magefile/mage/mg"
)

var (
	Docker = env.CrossPlatformBinary("docker")

	DefaultPlatforms []string
)

func setup() {
	mg.Deps(mg.F(tool.Exists, Docker), setDefaultBuildxPlatforms)
}

func setDefaultBuildxPlatforms() error {
	output, err := cmd.Run(Docker, "buildx", "inspect", "--bootstrap")
	if err != nil {
		return err
	}
	var setDefaultPlatforms bool
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		platformPrefix := "Platforms: "
		if !strings.HasPrefix(line, platformPrefix) {
			continue
		}
		DefaultPlatforms = strings.Split(strings.TrimPrefix(line, platformPrefix), ", ")
		setDefaultPlatforms = true
	}
	if !setDefaultPlatforms {
		return fmt.Errorf("unable to detect default platforms with docker buildx inspect --boostrap")
	}
	return nil
}
