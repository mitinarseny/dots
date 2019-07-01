package core

import (
	"errors"
	"os/exec"

	"gopkg.in/yaml.v3"
)

const (
	commandStringFormat = "> %s"
	commandOutputPrefix = "%s| "
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
	String      *string
	Description *string
	Commands    []*Command
	*exec.Cmd
}

type yamlCommandInline string

type yamlCommandExtended struct {
	Command     *yamlCommandInline
	Description *string
	Commands    []*Command
}

func (c *Command) UnmarshalYAML(value *yaml.Node) error {
	var auxInline yamlCommandInline
	if err := value.Decode(&auxInline); err == nil {
		c.String = (*string)(&auxInline)
		return nil
	}

	var auxExtended yamlCommandExtended
	if err := value.Decode(&auxExtended); err == nil {
		c.String = (*string)(auxExtended.Command)
		c.Description = auxExtended.Description
		c.Commands = auxExtended.Commands
		return nil
	}

	return errors.New("unable to parse command")
}

func (c *Command) Inspect() error {
	if c.String != nil {
		c.Cmd = exec.Command("/bin/sh", "-c", *c.String)
	}
	for _, subCmd := range c.Commands {
		if err := subCmd.Inspect(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) WithString(s string) *Command {
	c.String = &s
	c.Cmd = exec.Command("sh", "-c", *c.String)
	return c
}

func (c *Command) Execute() error {
	if c.Description != nil {
		logger.Printf("%s: ", *c.Description)
	}
	if c.Cmd != nil {
		if c.String != nil {
			logger.Printf(commandStringFormat, *c.String)
		}
		defer logger.SetPrefixf(commandOutputPrefix)()

		c.Cmd.Stdout = loggerWriter()

		if err := c.Cmd.Run(); err != nil {
			return err
		}
	}

	defer logger.SetPrefixf("  %s")()
	for _, subCmd := range c.Commands {
		if err := subCmd.Execute(); err != nil {
			return err
		}
	}
	return nil
}

func ExecuteCommands(cmds ...*Command) error {
	for _, c := range cmds {
		if err := c.Execute(); err != nil {
			return err
		}
	}
	return nil
}
