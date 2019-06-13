package core

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"sync"
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
	fmt.Println("Links:")
	if err := h.Link(); err != nil {
		return err
	}

	// TODO: other stages

	return nil
}

func (h *Host) Link() error {
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

	return WaitForPipeline(ctx, errCh)
}

func WaitForPipeline(ctx context.Context, errChs ...<-chan error) error {
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

func mergeErrors(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	// We must ensure that the output channel has the capacity to
	// hold as many errors
	// as there are error channels.
	// This will ensure that it never blocks, even
	// if WaitForPipeline returns early.
	out := make(chan error, len(cs))
	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls
	// wg.Done.
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	// Start a goroutine to close out once all the output goroutines
	// are done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

//func mergeToLinkChs(cs ...<-chan *ToLink) <-chan *ToLink {
//	var wg sync.WaitGroup
//	out := make(chan *ToLink)
//
//	// Start an output goroutine for each input channel in cs.  output
//	// copies values from c to out until c is closed, then calls wg.Done.
//	output := func(c <-chan *ToLink) {
//		for n := range c {
//			out <- n
//		}
//		wg.Done()
//	}
//	wg.Add(len(cs))
//	for _, c := range cs {
//		go output(c)
//	}
//
//	// Start a goroutine to close out once all the output goroutines are
//	// done.  This must start after the wg.Add call.
//	go func() {
//		wg.Wait()
//		close(out)
//	}()
//	return out
//}
