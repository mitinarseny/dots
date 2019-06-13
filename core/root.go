package core

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Config map[string]*Host

type yamlConfig map[string]*Host

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	var aux yamlConfig
	if err := value.Decode(&aux); err != nil {
		return err
	}
	*c = make(Config)
	for hostName, host := range aux {
		if host.Extends != nil {
			extendHost, ok := aux[host.Extends.Name]
			if !ok {
				return fmt.Errorf("unable to extend %q with %q", hostName, host.Name)
			}
			host.Extends = extendHost
		}
		host.Name = hostName
		(*c)[hostName] = host
	}

	return nil
}

func (c *Config) Inspect() error {
	for _, host := range *c {
		if err := host.Inspect(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) Up(hostName string) error {
	host, ok := (*c)[hostName]
	if !ok {
		return fmt.Errorf("there is no host %q", hostName)
	}

	if err := host.Up(); err != nil {
		return err
	}
	return nil
}
