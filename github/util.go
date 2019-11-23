package github

import (
	"github.com/hilalisadev/hub/git"
)

func IsHttpsProtocol() bool {
	httpProtocol, _ := git.Config("hub.protocol")
	return httpProtocol == "https"
}
