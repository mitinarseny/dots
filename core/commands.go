package core

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os/exec"
)

type Commands []Command

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
