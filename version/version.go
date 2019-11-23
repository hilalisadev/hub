package version

import (
	"fmt"

	"github.com/hilalisadev/hub/git"
)

var Version = "2.13.0"

func FullVersion() (string, error) {
	gitVersion, err := git.Version()
	if err != nil {
		gitVersion = "git version (unavailable)"
	}
	return fmt.Sprintf("%s\nhub version %s", gitVersion, Version), err
}
