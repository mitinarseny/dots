package config

type Config struct {
	Source string
	Host   `yaml:",inline"`
	Hosts  map[string]*Host
}

type Reviser interface {
	Revise(configDir string) error
}

func (c *Config) Revise(configDir string) error {
	links := c.Links
	for _, hv := range c.Hosts {
		links = append(links, hv.Links...)
	}
	return nil
}
