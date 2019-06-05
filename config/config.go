package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

type FilePath struct {
	Origin   string
	Absolute string
}

type Link struct {
	Source struct {
		Original string
		Relative string
	}
	Target struct {
		Original string
		Absolute string
	}
	Force bool
}

func (l *Link) Revise(configDir string) error {
	t, err := homedir.Expand(l.Target.Original)
	if err != nil {
		return err
	}
	t, err = filepath.Abs(t)
	if err != nil {
		return err
	}
	l.Target.Absolute = t

	s := filepath.Join(configDir, l.Source.Original)
	if _, err := os.Stat(s); err != nil {
		return err
	}
	l.Source.Relative = s

	return nil
}

type Variable struct {
	Name  string
	Value string
}

type DotfilesConfig struct {
	Source    string
	Variables map[string]string
	Links     []*Link
}

func (c *DotfilesConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var aux struct {
		Variables []struct {
			Name  string
			Value string
		}
		Links map[string]struct {
			Path  string
			Force bool
		}
	}

	if err := unmarshal(&aux); err != nil {
		return err
	}

	mapper := func(s string) string {
		return c.Variables[s]
	}

	for _, v := range aux.Variables {
		c.Variables[v.Name] = os.Expand(v.Value, mapper)
	}

	linksChan := make(chan *Link)
	for t, s := range aux.Links {
		go func(target, source string, force bool) {
			l := new(Link)
			l.Target.Original = os.Expand(target, mapper)
			l.Source.Original = os.Expand(source, mapper)
			l.Force = force
			linksChan <- l
		}(t, s.Path, s.Force)
	}

	for range aux.Links {
		c.Links = append(c.Links, <-linksChan)
	}

	return nil
}

func (c *DotfilesConfig) Revise() error {
	configDir := filepath.Dir(c.Source)

	ch := make(chan error)
	for _, ll := range c.Links {
		go func(l *Link) {
			if err := l.Revise(configDir); err != nil {
				ch <- fmt.Errorf("'%v' -> '%v': %v", l.Source.Original, l.Target.Original, err.Error())
				return
			}
			ch <- nil
		}(ll)
	}
	for range c.Links {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}
