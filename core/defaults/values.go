package defaults

import (
	"fmt"
)

type String string
type yamlString string

func (s *String) String() string {
	return fmt.Sprintf("-string %s", *s)
}

type Data string
type yamlData string

func (d *Data) String() string {
	return fmt.Sprintf("-data %s", *d)
}

type Integer int64

func (r *Integer) String() string {
	return fmt.Sprintf("-int %d", *r)
}

type Float float64

func (r *Float) String() string {
	return fmt.Sprintf("-float %f", *r)
}

type Boolean bool

func (r *Boolean) String() string {
	return fmt.Sprintf("-bool %t", *r)
}

type Date string

func (r *Date) String() string {
	return fmt.Sprintf("-string %s", *r)
}

type Array []string

func (r *Array) String() (s string) {
	s += "-array"
	for _, v := range *r {
		s += fmt.Sprintf(" %s", v)
	}
	return
}

type ArrayAdd []string

func (r *ArrayAdd) String() (s string) {
	s += "-array"
	for _, v := range *r {
		s += fmt.Sprintf(" %s", v)
	}
	return
}

type Dict map[string]string

func (r *Dict) String() (s string) {
	s += "-dict"
	for k, v := range *r {
		s += fmt.Sprintf(" %s %s", k, v)
	}
	return
}

type DictAdd map[string]string

func (r *DictAdd) String() (s string) {
	s += "-dict-add"
	for k, v := range *r {
		s += fmt.Sprintf(" %s %s", k, v)
	}
	return
}
