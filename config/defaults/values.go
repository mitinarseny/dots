package defaults

import (
	"fmt"
	"strconv"
)

type String string
type yamlString string

func (s *String) String() string {
	return fmt.Sprintf("-string %s", *s)
}

//func (s *String) UnmarshalYAML(value *yaml.Node) error {
//	var aux yamlString
//	if err := value.Decode(&aux); err != nil {
//		return err
//	}
//	*s = String(aux)
//	return nil
//}

func (s *String) Scan(i interface{}) error {
	switch v := i.(type) {
	case string:
		*s = String(v)
	default:
		return fmt.Errorf("cannot scan %T from %T", s, v)
	}
	return nil
}

type Data string
type yamlData string

func (d *Data) String() string {
	return fmt.Sprintf("-data %s", *d)
}

//func (d *Data) UnmarshalYAML(value *yaml.Node) error {
//	var aux yamlData
//	if err := value.Decode(&aux); err != nil {
//		return err
//	}
//	*d = Data(aux)
//	return nil
//}

func (d *Data) Scan(i interface{}) error {
	switch v := i.(type) {
	case string:
		*d = Data(v)
	default:
		return fmt.Errorf("cannot scan %T from %T", d, v)
	}
	return nil
}

type Integer int64

func (r *Integer) String() string {
	return fmt.Sprintf("-int %d", *r)
}

func (r *Integer) Scan(i interface{}) error {
	switch v := i.(type) {
	case string:
		vi, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		*r = Integer(vi)
	default:
		return fmt.Errorf("cannot scan %T from %T", r, v)
	}
	return nil
}

type Float float64

func (r *Float) String() string {
	return fmt.Sprintf("-float %f", *r)
}

func (r *Float) Scan(i interface{}) error {
	switch v := i.(type) {
	case string:
		vi, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		*r = Float(vi)
	default:
		return fmt.Errorf("cannot scan %T from %T", r, v)
	}
	return nil
}

type Boolean bool

func (r *Boolean) String() string {
	return fmt.Sprintf("-bool %t", *r)
}

func (r *Boolean) Scan(i interface{}) error {
	switch v := i.(type) {
	case string:
		vi, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		*r = Boolean(vi)
	default:
		return fmt.Errorf("cannot scan %T from %T", r, v)
	}
	return nil
}

type Date string

func (r *Date) String() string {
	return fmt.Sprintf("-string %s", *r)
}

func (r *Date) Scan(i interface{}) error {
	switch v := i.(type) {
	case string:
		*r = Date(v)
	default:
		return fmt.Errorf("cannot scan %T from %T", r, v)
	}
	return nil
}

type Array []string

func (r *Array) String() (s string) {
	s += "-array"
	for _, v := range *r {
		s += fmt.Sprintf(" %s", v)
	}
	return
}

func (r *Array) Scan(i interface{}) error {
	switch v := i.(type) {
	case []string:
		if v == nil {
			return fmt.Errorf("cannot scan %T from empty array", r)
		}
		*r = v
	default:
		return fmt.Errorf("cannot scan %T from %T", r, v)
	}
	return nil
}

type ArrayAdd []string

func (r *ArrayAdd) String() (s string) {
	s += "-array"
	for _, v := range *r {
		s += fmt.Sprintf(" %s", v)
	}
	return
}

func (r *ArrayAdd) Scan(i interface{}) error {
	switch v := i.(type) {
	case []string:
		if v == nil {
			return fmt.Errorf("cannot scan %T from empty array", r)
		}
		*r = v
	default:
		return fmt.Errorf("cannot scan %T from %T", r, v)
	}
	return nil
}

type Dict map[string]string

func (r *Dict) String() (s string) {
	s += "-dict"
	for k, v := range *r {
		s += fmt.Sprintf(" %s %s", k, v)
	}
	return
}

func (r *Dict) Scan(i interface{}) error {
	switch v := i.(type) {
	case map[string]string:
		if v == nil {
			return fmt.Errorf("cannot scan %T from empty dict", r)
		}
		*r = v
	default:
		return fmt.Errorf("cannot scan %T from %T", r, v)
	}
	return nil
}

type DictAdd map[string]string

func (r *DictAdd) String() (s string) {
	s += "-dict-add"
	for k, v := range *r {
		s += fmt.Sprintf(" %s %s", k, v)
	}
	return
}

func (r *DictAdd) Scan(i interface{}) error {
	switch v := i.(type) {
	case map[string]string:
		if v == nil {
			return fmt.Errorf("cannot scan %T from empty map", r)
		}
		*r = v
	default:
		return fmt.Errorf("cannot scan %T from %T", r, v)
	}
	return nil
}
