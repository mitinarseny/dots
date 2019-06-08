package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"os/exec"
	"sync"
)

type Commands []*Command

type Command string

func (c *Command) UnmarshalYAML(value *yaml.Node) error {
	var cmd string

	if err := value.Decode(&cmd); err != nil {
		return err
	}

	*c = Command(cmd)

	return nil
}

func (c *Command) Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	cmd := exec.Command("sh", "-c", string(*c))

	var wg sync.WaitGroup

	var errStdin, errStdout, errStderr error

	//cmdStdin, err := cmd.StdinPipe()
	//if err != nil {
	//	return err
	//}
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	_ , errStdin = io.Copy(cmdStdin, stdin)
	//}()

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, errStdout = io.Copy(stdout, cmdStdout)
	}()

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, errStderr = io.Copy(stderr, cmdStderr)
	}()

	if err := cmd.Start(); err != nil {
		return err
	}
	wg.Wait()

	for _, err := range []error{errStdin, errStdout, errStderr}{
		if err != nil {
			return err
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
