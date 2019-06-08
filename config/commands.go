package config

import (
	"gopkg.in/yaml.v3"
	"os/exec"
)

type Commands []Command

type Command struct {
	String string
	*exec.Cmd
}

func (c *Command) WithString(s string) *Command {
	c.String = s
	c.Cmd = exec.Command("sh", "-c", c.String)
	return c
}

type yamlCommand string

func (c *Command) UnmarshalYAML(value *yaml.Node) error {
	var aux yamlCommand

	if err := value.Decode(&aux); err != nil {
		return err
	}

	c.WithString(string(aux))

	return nil
}
