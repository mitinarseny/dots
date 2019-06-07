package config

import (
	"github.com/mitchellh/go-homedir"
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
		ll.Target = &FilePath{
			Original:target,
		}
		*l = append(*l, &ll)
	}
	return nil
}

type Link struct {
	Source *FilePath
	Target *FilePath
	Force  bool
}

type yamlLink struct {
	Path *FilePath
	Force  bool
}

func (l *Link) UnmarshalYAML(value *yaml.Node) error {
	var aux yamlLink
	if err := value.Decode(&aux); err != nil {
		return err
	}
	l.Source = aux.Path
	l.Force = aux.Force
	return nil
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
