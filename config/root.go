package config

type Config struct {
	Host   `yaml:",inline"`
	Hosts  map[string]*Host
}

type Reviser interface {
	Revise(configDir string) error
}

//func (c *Config) Revise(configDir string, hostName *string) error {
//	if err := c.Host.Revise(configDir); err != nil {
//		return err
//	}
//	if hostName != nil {
//		h, exists := c.Hosts[*hostName]
//		if !exists {
//			return fmt.Errorf("there is no host '%s'", *hostName)
//		}
//		if err := h.Revise(configDir); err != nil {
//			return err
//		}
//	}
//	return nil
//}
