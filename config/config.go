package config

type Link struct {
	Source string
	Target string
	Force  bool
}

type Variable struct {
	Name  string
	Value string
}

type Config struct {
	Variables []Variable
	Links     []Link
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Variables map[string]string
		Links     map[string]struct {
			Path  string
			Force bool
		}
	}

	if err := unmarshal(&aux); err != nil {
		return err
	}

	for n, v := range aux.Variables {
		c.Variables = append(c.Variables, Variable{
			Name:  n,
			Value: v,
		})
	}

	for t, s := range aux.Links {
		c.Links = append(c.Links, Link{
			Target: t,
			Source: s.Path,
			Force:  s.Force,
		})
	}

	return nil
}
