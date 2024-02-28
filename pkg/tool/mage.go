package tool

import (
	"os/exec"
)

// Exists checks if a particular tool exists on the host machine or returns an error
func Exists(tool string) error {
	_, err := exec.LookPath(tool)
	return err
}
