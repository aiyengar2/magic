package docker

import (
	"fmt"
	"os"

	"github.com/aiyengar2/magic/pkg/utils/cmd"
	"github.com/magefile/mage/mg"
)

func BuildxCreate() error {
	mg.Deps(setup)
	return cmd.RunV(Docker, "buildx", "create", "--use", "--name", "magic-multi-platform-builder", "--node", "magic")
}

func BuildxRemove() error {
	if len(os.Getenv("MAGIC_BUILDX_KEEP_ALIVE")) > 0 {
		fmt.Printf("Warning: skipping removal of docker buildx runner to speed up builds, please run 'docker buildx rm magic-multi-platform-builder' manually to clean up.\n")
		return nil
	}
	mg.Deps(setup)
	return cmd.RunV(Docker, "buildx", "rm", "magic-multi-platform-builder")
}
