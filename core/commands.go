package core

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os/exec"
)

const (
	commandOutputPrefix = "  |  "
)

type Commands []*Command

func (c *Commands) Inspect() error {
	for _, cmd := range *c {
		if err := cmd.Inspect(); err != nil {
			return err
		}
	}
	return nil
}

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
		c.String = string(auxInline)
		return nil
	}

	var auxExtended yamlCommandExtended
	if err := value.Decode(&auxExtended); err == nil {
		c.String = string(auxExtended.Command)
		c.Description = auxExtended.Description
		return nil
	}

	return errors.New("unable to parse command")
}

func (c *Command) Inspect() error {
	c.Cmd = exec.Command("sh", "-c", c.String)
	return nil
}

func (c *Command) WithString(s string) *Command {
	c.String = s
	c.Cmd = exec.Command("sh", "-c", c.String)
	return c
}

func (c *Command) Execute() error {
	defer logger.SetPrefixf("%s | ")()

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