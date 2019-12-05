package commands

import (
	"github.com/hilalisadev/hub/ui"
	"github.com/hilalisadev/hub/utils"
	"github.com/hilalisadev/hub/version"
)

var cmdVersion = &Command{
	Run:          runVersion,
	Usage:        "version",
	Long:         "Shows git version and hub client version.",
	GitExtension: true,
}

func init() {
	CmdRunner.Use(cmdVersion, "--version")
}

func runVersion(cmd *Command, args *Args) []byte {
	output, err := version.FullVersion()
	if output != "" {
		ui.Println(output)
	}
	utils.Check(err)
	args.NoForward()
	return nil
}
