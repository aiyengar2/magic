package git

import (
	"bytes"
	"strings"

	"github.com/aiyengar2/magic/pkg/utils/cmd"
)

func Commit() string {
	// Does not depend on setup since we don't return an error if the command fails
	// mg.Deps(setup)

	outputBuf, errBuf := &bytes.Buffer{}, &bytes.Buffer{}
	err := cmd.Exec(nil, nil, outputBuf, errBuf, Git, "rev-parse", "--short", "HEAD")
	output := strings.TrimSpace(outputBuf.String())

	if err != nil {
		cmd.PrintWarning("could not get git commit: %s\n%s", err.Error(), errBuf.String())
		return "HEAD"
	}
	return output
}

func Dirty() bool {
	// Does not depend on setup since we don't return an error if the command fails
	// mg.Deps(setup)

	outputBuf, errBuf := &bytes.Buffer{}, &bytes.Buffer{}
	err := cmd.Exec(nil, nil, outputBuf, errBuf, Git, "status", "--porcelain", "--untracked-files=no")
	output := strings.TrimSpace(outputBuf.String())

	if err != nil {
		cmd.PrintWarning("could not get git status: %s\n%s", err.Error(), errBuf.String())
		return false
	}
	return len(output) > 0
}

func Tag() string {
	// Does not depend on setup since we don't return an error if the command fails
	// mg.Deps(setup)

	outputBuf, errBuf := &bytes.Buffer{}, &bytes.Buffer{}
	err := cmd.Exec(nil, nil, outputBuf, errBuf, Git, "tag", "-l", "--contains", "HEAD")
	output := strings.TrimSpace(outputBuf.String())

	if err != nil {
		// if we can't get the tag, return nothing
		cmd.PrintWarning("could not get git tag: %s\n%s", err.Error(), errBuf.String())
		return ""
	}
	if len(output) == 0 {
		return ""
	}
	return strings.SplitN(output, "\n", 2)[0]
}
