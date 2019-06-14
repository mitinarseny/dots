package core

import (
	"context"
	"gopkg.in/yaml.v3"
	"sync"
)

const (
	stagePrefix = "  |  "
)

type Host struct {
	Name      string
	Extends   *Host
	Variables *Variables
	Links     *Links
	Commands  *Commands
	//Defaults  defaults.Defaults
}

type yamlHost struct {
	Extends   *string
	Variables *Variables
	Links     *Links
	Commands  *Commands
	//Defaults  defaults.Defaults
}

func (h *Host) UnmarshalYAML(value *yaml.Node) error {
	var aux yamlHost
	if err := value.Decode(&aux); err != nil {
		return err
	}
	if aux.Extends != nil {
		h.Extends = &Host{Name: *aux.Extends}
	}
	h.Variables = aux.Variables
	h.Links = aux.Links
	h.Commands = aux.Commands
	//h.Defaults = aux.Defaults

	return nil
}

func (h *Host) Inspect() error {
	if h.Links != nil {
		if err := h.Links.Inspect(); err != nil {
			return err
		}
	}

	// TODO: inspect others

	return nil
}

func (h *Host) Up() error {
	if err := h.SetVariables(); err != nil {
		return err
	}
	if err := h.CreateLinks(); err != nil {
		return err
	}
	if err := h.ExecuteCommands(); err != nil {
		return err
	}

	// TODO: other stages

	return nil
}

func (h *Host) SetVariables() error {
	logger.Println("Variables:")
	defer logger.SetPrefix(logger.Prefix())
	logger.SetPrefix(logger.Prefix() + stagePrefix)

	vars, err := h.CollectVariables()
	if err != nil {
		return err
	}

	return SetVariables(vars...)
}

func (h *Host) CollectVariables() ([]*Variable, error) {
	var vars []*Variable
	if h.Extends != nil && h.Extends.Variables != nil {
		hostVars, err := h.Extends.CollectVariables()
		if err != nil {
			return nil, err
		}
		vars = append(vars, hostVars...)
	}
	if h.Variables != nil {
		hostVars, err := h.Variables.GenVariables()
		if err != nil {
			return nil, err
		}
		vars = append(vars, hostVars...)
	}
	return vars, nil
}

func (h *Host) CreateLinks() error {
	logger.Println("Links:")
	defer logger.SetPrefix(logger.Prefix())
	logger.SetPrefix(logger.Prefix() + stagePrefix)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	toLinkCh, err := h.GenLinks(ctx)
	if err != nil {
		return err
	}

	errCh, err := Linker(ctx, toLinkCh)
	if err != nil {
		return err
	}

	return waitForPipeline(ctx, errCh)
}

func (h *Host) GenLinks(ctx context.Context) (<-chan *ToLink, error) {
	var toLinkChs []<-chan *ToLink

	if h.Extends != nil && h.Extends.Links != nil {
		toLinkCh, err := h.Extends.GenLinks(ctx)
		if err != nil {
			return nil, err
		}
		toLinkChs = append(toLinkChs, toLinkCh)
	}
	if h.Links != nil {
		for _, l := range *h.Links {
			toLinkCh, err := l.GenLinks(ctx)
			if err != nil {
				return nil, err
			}
			toLinkChs = append(toLinkChs, toLinkCh)
		}
	}
	return mergeToLinkChs(ctx, toLinkChs...), nil
}

func (h *Host) ExecuteCommands() error {
	logger.Println("Commands:")
	defer logger.SetPrefix(logger.Prefix())
	logger.SetPrefix(logger.Prefix() + stagePrefix)

	cmds, err := h.CollectCommands()
	if err != nil {
		return err
	}
	return ExecuteCommands(cmds...)
}

func (h *Host) CollectCommands() ([]*Command, error) {
	var cmds []*Command
	if h.Extends != nil && h.Extends.Commands != nil {
		hostCmds, err := h.Extends.CollectCommands()
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, hostCmds...)
	}
	if h.Commands != nil {
		hostCmds, err := h.Commands.CollectCommands()
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, hostCmds...)
	}
	return cmds, nil
}

func waitForPipeline(ctx context.Context, errChs ...<-chan error) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	errc := mergeErrorChs(ctx, errChs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

func mergeErrorChs(ctx context.Context, cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)

	output := func(c <-chan error) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		defer close(out)

		wg.Wait()
	}()
	return out
}

func mergeToLinkChs(ctx context.Context, cs ...<-chan *ToLink) <-chan *ToLink {
	var wg sync.WaitGroup
	out := make(chan *ToLink)

	output := func(c <-chan *ToLink) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		defer close(out)

		wg.Wait()
	}()
	return out
}
