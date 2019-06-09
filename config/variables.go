package config

import (
	"fmt"
	"go.uber.org/zap/buffer"
	"gopkg.in/yaml.v3"
	"strings"
)

type Variables []map[string]*Variable

type Variable struct {
	Value   *string
	Command *Command
}

type yamlVariableInline *string

type yamlVariableExtended struct {
	Value   *string
	Command *Command
}

func (v *Variable) UnmarshalYAML(value *yaml.Node) error {
	var auxInline yamlVariableInline
	if err := value.Decode(&auxInline); err == nil {
		v.Value = auxInline
		return nil
	}

	var auxExtended yamlVariableExtended
	if err := value.Decode(&auxExtended); err == nil {
		if auxExtended.Value != nil && auxExtended.Command != nil {
			return fmt.Errorf("variable should be <string>, or { value: <string> }, or { command: <command> }")
		}
		v.Value = auxExtended.Value
		v.Command = auxExtended.Command
		return nil
	}
	return fmt.Errorf("variable should be <string>, or { value: <string> }, or { command: <command> }")
}

func (v *Variable) FromCommand() error {
	var out buffer.Buffer
	v.Command.Stdout = &out
	if err := v.Command.Run(); err != nil {
		return err
	}
	varVal := strings.TrimSpace(out.String())
	v.Value = &varVal
	return nil
}
