// +build go1.8

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/hilalisadev/hub/commands"
	"github.com/hilalisadev/hub/github"
	"github.com/hilalisadev/hub/ui"
)

func main() {
	defer github.CaptureCrash()
	err, toto := commands.CmdRunner.Execute(os.Args)
	fmt.Println("dslkjflkj", toto)
	exitCode := handleError(err)
	os.Exit(exitCode)
}

func handleError(err error) int {
	if err == nil {
		return 0
	}

	switch e := err.(type) {
	case *exec.ExitError:
		if status, ok := e.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		} else {
			return 1
		}
	case *commands.ErrHelp:
		ui.Println(err)
		return 0
	default:
		if errString := err.Error(); errString != "" {
			ui.Errorln(err)
		}
		return 1
	}
}
