package config

import "gopkg.in/yaml.v3"

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