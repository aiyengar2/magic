package version

import (
	"fmt"

	"github.com/aiyengar2/magic/pkg/git"
)

var (
	Dirty   bool
	Commit  = "HEAD"
	GitTag  string
	Version = "v0.0.0-dev"
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

}

func GetVersion() string {
	return fmt.Sprintf("%s (%s)", Version, Commit)
}
