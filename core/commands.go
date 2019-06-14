package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os/exec"
	"strings"
)

const (
	commandOutputPrefix = "  |  "
)

type Commands []*Command

func (c *Commands) CollectCommands() ([]*Command, error) {
	return *c, nil
}

type Command struct {
	String      string
	Description *string
	*exec.Cmd
}

type yamlCommandInline string

type yamlCommandExtended struct {
	Command     yamlCommandInline
	Description *string
}

func (c *Command) UnmarshalYAML(value *yaml.Node) error {
	var auxInline yamlCommandInline
	if err := value.Decode(&auxInline); err == nil {
		c.WithString(string(auxInline))
		return nil
	}

	var auxExtended yamlCommandExtended
	if err := value.Decode(&auxExtended); err == nil {
		c.WithString(string(auxExtended.Command)).Description = auxExtended.Description
		return nil
	}

	return errors.New("unable to parse command")
}

func (c *Command) WithString(s string) *Command {
	c.String = s
	c.Cmd = exec.Command("sh", "-c", c.String)
	return c
}

func (c *Command) Execute() error {
	defer logger.SetPrefix(logger.Prefix())

	logger.SetPrefix(strings.Repeat(" ", len(logger.Prefix())) + commandOutputPrefix)

	c.Cmd.Stdout = loggerWriter()

	if err := c.Cmd.Run(); err != nil {
		return err
	}
	return nil
}

func ExecuteCommands(cmds ...*Command) error {
	for _, c := range cmds {
		var buff bytes.Buffer
		if c.Description != nil {
			buff.WriteString(fmt.Sprintf("%s: ", *c.Description))
		}
		logger.Println(buff.String() + c.String)
		if err := c.Execute(); err != nil {
			return err
		}
	}
	return nil
}

func prefixedWriter(w io.Writer, prefix string) io.Writer {
	pipeReader, pipeWriter := io.Pipe()

	scanner := bufio.NewScanner(pipeReader)
	scanner.Split(bufio.ScanLines)

	go func() {
		for scanner.Scan() {
			_, _ = fmt.Fprint(w, prefix)
			_, _ = w.Write(scanner.Bytes())
			_, _ = fmt.Fprint(w, '\n')
		}
	}()
	return pipeWriter
}
