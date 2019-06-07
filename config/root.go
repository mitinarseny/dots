package config

import (
	"fmt"
	"path/filepath"
)

type Config struct {
	Source string
	Host   `yaml:",inline"`
	Hosts  map[string]*Host
}

//type root struct {
//	yamlHost
//	Hosts map[string]yamlHost
//}

//func (c *Config) UnmarshalYAML(value *yaml.Node) error {
//
//	var aux root
//
//	if err := value.Decode(&aux); err != nil {
//		return err
//	}
//	for hostName, hostValue := range aux.Hosts {
//
//	}
//
//	return nil
//}

//func scanDefaults(items map[string]map[string]interface{}) (domains map[string]defaults.Domain, err error) {
//	domains = make(map[string]defaults.Domain)
//	for appName, keys := range items {
//		domains[appName] = make(map[string]*defaults.Key)
//		for key, value := range keys {
//			domains[appName][key] = new(defaults.Key)
//			if err := domains[appName][key].Scan(value); err != nil {
//				return nil, err
//			}
//		}
//	}
//	return
//}

func (c *Config) Revise() error {
	configDir := filepath.Dir(c.Source)

	ch := make(chan error)
	for _, l := range c.Links {
		go func(l *Link) {
			if err := l.Revise(configDir); err != nil {
				ch <- fmt.Errorf("'%v' -> '%v': %v", l.Source.Original, l.Target.Original, err.Error())
				return
			}
			ch <- nil
		}(l)
	}
	for range c.Links {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}
