package config

import (
	"github.com/mitinarseny/dots/config/defaults"
)

type Host struct {
	Variables Variables
	Links     Links
	Commands  Commands
	Defaults  defaults.Defaults
}

func (h *Host) Revise(configDir string) error {
	for _, l := range h.Links {
		if err := l.Revise(configDir); err != nil {
			return err
		}
	}
	return nil
}