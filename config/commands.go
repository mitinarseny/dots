package config

import (
	"io"
	"os/exec"
)

type Commands []*Command

type Command string

func (c *Command) Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	cmd := exec.Command("sh", "-c", string(*c))

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
