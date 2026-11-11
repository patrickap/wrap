package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/patrickap/wrap/cmd"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		var exitError *exec.ExitError

		if errors.As(err, &exitError) {
			os.Exit(exitError.ExitCode())
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
