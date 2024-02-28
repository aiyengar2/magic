package version

import (
	"fmt"
	"os"

	"github.com/aiyengar2/magic/pkg/git"
)

var (
	Dirty   bool
	Commit  = "HEAD"
	GitTag  string
	Version = "v0.0.0-dev"

	Org = "arvindiyengar" // TODO: modify this when finalizing the repository
	Tag = "v0.0.0-dev"
)

func init() {
	Dirty = git.Dirty()
	Commit = git.Commit()
	GitTag = git.Tag()

	if !Dirty && len(GitTag) > 0 {
		Version = GitTag
	} else {
		if Dirty {
			Version = fmt.Sprintf("%s-dirty", Commit)
		} else {
			Version = Commit
		}
	}

	// Get org
	if org := os.Getenv("ORG"); len(org) > 0 {
		Org = org
	}

	// Get tag
	Tag = os.Getenv("TAG")
	if len(Tag) == 0 {
		Tag = Version
	}
}

func GetVersion() string {
	return fmt.Sprintf("%s (%s)", Version, Commit)
}

func GetTag(repo string) string {
	return fmt.Sprintf("%s/%s:%s", Org, repo, Tag)
}
