package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const (
	Omitted  = 0
	Created  = 1
	Replaced = 2
)

type Links []*Link

type yamlLinks map[string]*Link

func (l *Links) UnmarshalYAML(value *yaml.Node) error {
	var aux yamlLinks
	if err := value.Decode(&aux); err != nil {
		return err
	}
	for target, link := range aux {
		ll := *link
		ll.Target.Original = target
		*l = append(*l, &ll)
	}
	return nil
}

type Link struct {
	Source FilePath
	Target FilePath
	Force  bool
}

type yamlLinkInline *FilePath

type yamlLinkExtended struct {
	Path  FilePath
	Force bool
}

func (l *Link) UnmarshalYAML(value *yaml.Node) error {
	var auxInline yamlLinkInline
	if err := value.Decode(&auxInline); err == nil {
		l.Source = *auxInline
		l.Force = false
		return nil
	}
	var aux yamlLinkExtended
	if err := value.Decode(&aux); err == nil {
		l.Source = aux.Path
		l.Force = aux.Force
		return nil
	}
	return fmt.Errorf("link should be <string>, or { path: <string>, force: <bool> }")
}

type FilePath struct {
	Original string
	Absolute string
}

type yamlFilePath string

func (f *FilePath) UnmarshalYAML(value *yaml.Node) error {
	var aux yamlFilePath
	if err := value.Decode(&aux); err != nil {
		return err
	}
	f.Original = string(aux)
	return nil
}

func (l *Link) Link() (status int, err error) {
	var existed bool
	switch _, err := os.Lstat(l.Target.Absolute); {
	case err == nil:
		existed = true
		if !l.Force {
			return Omitted, nil
		}
		if err := os.Remove(l.Target.Absolute); err != nil {
			return -1, err
		}
	case os.IsNotExist(err):
	default:
		return -1, err
	}
	if err := os.Symlink(l.Source.Absolute, l.Target.Absolute); err != nil {
		return -1, err
	}
	if existed {
		return Replaced, nil
	}
	return Created, nil
}
func (l *Link) Revise(configDir string) error {
	t, err := expandTilde(l.Target.Original)
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

func expandTilde(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(hd, path[1:]), nil
}
