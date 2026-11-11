package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Run(program string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	command := exec.Command(program, args...)

	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("run %q: %w", program, err)
	}

	return nil
}

func RunStdio(program string, args []string) error {
	return Run(program, args, os.Stdin, os.Stdout, os.Stderr)
}
