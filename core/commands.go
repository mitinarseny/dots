package core

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
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

func ExecuteCommands(cmds... *Command) error {
	for _, c := range cmds {
		var buff bytes.Buffer
		if c.Description != nil {
			buff.WriteString(fmt.Sprintf("%s: ", *c.Description))
		}
		buff.WriteString(c.String)
		fmt.Println(buff.String())
		c.Cmd.Stdout = os.Stdout
		if err := c.Cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}