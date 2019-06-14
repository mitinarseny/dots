package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Links map[string]*Link

type yamlLinks map[string]*Link

func (l *Links) UnmarshalYAML(value *yaml.Node) error {
	var aux yamlLinks
	if err := value.Decode(&aux); err != nil {
		return err
	}
	*l = make(Links, len(aux))
	for target, link := range aux {
		if link == nil {
			(*l)[target] = &Link{
				Target: target,
				Force:  false,
			}
			continue
		}
		link.Target = target
		(*l)[target] = link
	}
	return nil
}

func (l *Links) Inspect() error {
	for _, link := range *l {
		if err := link.Inspect(); err != nil {
			return err
		}
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
	Path  *string
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
		if aux.Path != nil {
			l.Path = *aux.Path
		}
		l.Force = aux.Force
		return nil
	}
	return fmt.Errorf("link should be <string>, or { path: <string>, force: <bool> }")
}

func (l *Link) Inspect() error {
	if isGlob(l.Path) && !strings.HasSuffix(l.Target, string(os.PathSeparator)) {
		return errors.New("if path is a directory or glob, Target must end with '/'")
	}
	if l.Path == "" || strings.HasSuffix(l.Path, string(os.PathSeparator)) {
		if isGlob(l.Target) {
			return fmt.Errorf("cannot deduce path from target %q", l.Target)
		}
		l.Path = filepath.Base(l.Target)
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
				Target: path.Join(absTargetExpanded, strings.TrimPrefix(source, basePath)),
				Force:  l.Force,
			}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

func Linker(ctx context.Context, in <-chan *ToLink) (<-chan error, error) {
	errCh := make(chan error)
	go func() {
		defer close(errCh)

		for tl := range in {
			err := tl.Link()
			select {
			case errCh <- err:
			case <-ctx.Done():
				return
			}
		}
	}()
	return errCh, nil
}

type ToLink struct {
	Source string
	Target string
	Force  bool
}

func (t *ToLink) Link() error {
	pathDir := path.Dir(t.Target)
	if err := os.MkdirAll(pathDir, 0755); err != nil {
		logger.Printf("failed to create path: %s", pathDir)
		return err
	}
	var buff bytes.Buffer
	defer func() {
		logger.Println(buff.String())
	}()

	buff.WriteString(fmt.Sprintf("%s <- %s: ", t.Target, t.Source))
	var existed bool
	switch _, err := os.Lstat(t.Target); {
	case err == nil:
		existed = true
		if !t.Force {
			buff.WriteString("omitted")
			return nil
		}
		// TODO: backup
		if err := os.RemoveAll(t.Target); err != nil {
			buff.WriteString("failed to remove existing")
			return err
		}
	case os.IsNotExist(err):
	default:
		buff.WriteString(err.Error())
		return err
	}

	if err := os.Symlink(t.Source, t.Target); err != nil {
		buff.WriteString(err.Error())
		return err
	}

	if existed {
		buff.WriteString("replaced")
		return nil
	}
	buff.WriteString("created")
	return nil
}

func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
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
