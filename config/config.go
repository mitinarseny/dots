package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"sync"
)

type DotfilesConfig struct {
	Source    string
	Variables map[string]string
	Links     []*Link
	Commands  []string
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
		Commands []string
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

	ch := make(chan *Link)
	var wg sync.WaitGroup
	wg.Add(len(aux.Links))
	for t, s := range aux.Links {
		go func(target, source string, force bool) {
			defer wg.Done()

			l := new(Link)
			l.Target.Original = os.Expand(target, mapper)
			l.Source.Original = os.Expand(source, mapper)
			l.Force = force

			ch <- l
		}(t, s.Path, s.Force)
	}
	for range aux.Links {
		c.Links = append(c.Links, <-ch)
	}
	c.Commands = aux.Commands
	return nil
}

func (c *DotfilesConfig) Revise() error {
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

type Link struct {
	Source FilePath
	Target FilePath
	Force  bool
}

type FilePath struct {
	Original string
	Absolute string
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

	s, err := filepath.Abs(filepath.Join(configDir, l.Source.Original))
	if err != nil {
		return err
	}
	if _, err := os.Stat(s); err != nil {
		return err
	}
	l.Source.Absolute = s


	go func() {

	}()


	return nil
}
