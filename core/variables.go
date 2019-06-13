package core

import (
	"fmt"
	"go.uber.org/zap/buffer"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Variables []VariableStage

func (v *Variables) Set() error {
	for _, vs := range *v {
		if err := vs.Set(); err != nil {
			return err
		}
	}
	return nil
}

type VariableStage map[string]*Variable

type yamlVariableStage map[string]*Variable

func (vs *VariableStage) UnmarshalYAML(value *yaml.Node) error {
	var aux yamlVariableStage
	if err := value.Decode(&aux); err != nil {
		return err
	}

	*vs = make(VariableStage, len(aux))
	for varName, variable := range aux {
		variable.Name = varName
		(*vs)[varName] = variable
	}

	return nil
}

func (vs *VariableStage) Set() error {
	for _, v := range *vs {
		if err := v.Set(); err != nil {
			return err
		}
	}
	return nil
}

type Variable struct {
	Name    string
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
func (v *Variable) Set() error {
	if v.Command != nil {
		var out buffer.Buffer
		v.Command.Stdout = &out
		if err := v.Command.Run(); err != nil {
			return err
		}
		varVal := strings.TrimSpace(out.String())
		v.Value = &varVal
	}
	if err := os.Setenv(v.Name, *v.Value); err != nil {
		return err
	}
	return nil
}