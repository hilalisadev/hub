// +build !windows

package github

import (
	"os"

	"github.com/hilalisadev/hub/cmd"
)

func setConsole(cmd *cmd.Cmd) {

	stdin, err := os.OpenFile("/dev/tty", os.O_RDONLY, 0660)
	if err == nil {
		cmd.Stdin = stdin
	}
}
