package core

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	Omitted  = LinkStatus(0)
	Created  = LinkStatus(1)
	Replaced = LinkStatus(2)
)

type LinkStatus int

type Links []*Link

type yamlLinks map[string]*Link

func (l *Links) UnmarshalYAML(value *yaml.Node) error {
	var aux yamlLinks
	if err := value.Decode(&aux); err != nil {
		return err
	}
	//*l = make(Links, len(aux))
	for target, link := range aux {
		link.Target = target
		*l = append(*l, link)
	}
	return nil
}

type Link struct {
	Target string
	Path   string
	Force  bool
}

type yamlLinkInline string

type yamlLinkExtended struct {
	Path  string
	Force bool
}

func (l *Link) UnmarshalYAML(value *yaml.Node) error {
	var auxInline yamlLinkInline
	if err := value.Decode(&auxInline); err == nil {
		l.Path = string(auxInline)
		l.Force = false
		return nil
	}
	var aux yamlLinkExtended
	if err := value.Decode(&aux); err == nil {
		l.Path = aux.Path
		l.Force = aux.Force
		return nil
	}
	return fmt.Errorf("link should be <string>, or { path: <string>, force: <bool> }")
}

func (l *Link) Inspect() error {
	if l.Path == "" {
		return errors.New("path cannot be empty")
	}
	if isGlob(l.Path) && !strings.HasSuffix(l.Target, string(os.PathSeparator)) {
		return errors.New("if path is a directory or glob, Target has to end with '/'")
	}
	return nil
}

func (l *Link) GenLinks(ctx context.Context) (<-chan *ToLink, error) {
	targetExpanded, err := expandTilde(l.Target)
	if err != nil {
		return nil, err
	}
	absTargetExpanded, err := filepath.Abs(targetExpanded)
	if err != nil {
		return nil, err
	}
	absSource, err := filepath.Abs(l.Path)
	if err != nil {
		return nil, err
	}

	out := make(chan *ToLink)

	if strings.HasSuffix(absSource, string(os.PathSeparator)) {
		absSource = absSource + "*"
	}
	if !isGlob(absSource) {
		if strings.HasSuffix(l.Target, string(os.PathSeparator)) {
			absTargetExpanded = path.Join(absTargetExpanded, filepath.Base(l.Path))
		}
		go func() {
			defer close(out)

			select {
			case out <- &ToLink{
				Source: absSource,
				Target: absTargetExpanded,
				Force:  l.Force,
			}:
			case <-ctx.Done():
				return
			}
		}()
		return out, nil
	}

	if !strings.HasSuffix(l.Target, string(os.PathSeparator)) {
		return nil, errors.New("cannot link glob path to single file")
	}

	basePath := latestNoGlob(absSource)
	sources, err := filepath.Glob(absSource)
	if err != nil {
		return nil, err
	}
	go func() {
		defer close(out)

		for _, source := range sources {
			select {
			case out <- &ToLink{
				Source: source,
				Target: path.Join(l.Target, strings.TrimPrefix(source, basePath)),
				Force:  l.Force,
			}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

//func Linker(ctx context.Context, in <-chan *ToLink) error {
//	for {
//		select {
//		case l := <-in:
//			if err := l.Link(); err != nil {
//				return err
//			}
//			fmt.Println("")
//		case <-ctx.Done():
//			return nil
//		}
//	}
//	out := make(chan LinkStatus)
//	errCh := make(chan error, 1)
//
//	go func() {
//		defer close(out)
//		defer close(errCh)
//		for toLink := range in {
//			st, err := toLink.Link()
//			if err != nil {
//				errCh <- err
//				return
//			}
//			select {
//			case out <- st:
//			case <-ctx.Done():
//				return
//			}
//		}
//	}()
//	return out, errCh, nil
//}

type ToLink struct {
	Source string
	Target string
	Force  bool
}

func (t *ToLink) Link() error {
	if err := os.MkdirAll(path.Dir(t.Target), 0755); err != nil {
		return err
	}

	var existed bool
	switch _, err := os.Lstat(t.Target); {
	case err == nil:
		existed = true
		if !t.Force {
			return nil
		}
		// TODO: backup
		if err := os.Remove(t.Target); err != nil {
			return err
		}
	case os.IsNotExist(err):
	default:
		return err
	}

	if err := os.Symlink(t.Source, t.Target); err != nil {
		return err
	}

	if existed {
		return nil
	}
	return nil
}

func expandTilde(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		return path, nil
	}
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(hd, path[1:]), nil
}

func latestNoGlob(path string) string {
	if !strings.ContainsAny(path, "*?[^]") {
		return path
	}
	for dir := filepath.Dir(path); ; dir = filepath.Dir(dir) {
		if !strings.ContainsAny(dir, "*?[^]") {
			return dir
		}
	}
}

func isGlob(path string) bool {
	return strings.ContainsAny(path, "*?[^]")
}
