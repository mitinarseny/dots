package defaults

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

const (
	StringType   = "string"
	DataType     = "data"
	IntegerType  = "int"
	FloatType    = "float"
	BooleanType  = "bool"
	DateType     = "date"
	ArrayType    = "array"
	ArrayAddType = "array-add"
	DictType     = "dict"
	DictAddType  = "dict-add"
)

type Key struct {
	Type  string
	Value interface {
		fmt.Stringer
	}
}

type yamlKeyInline string

type yamlKeyExtendedCommon struct {
	Type  string
}

type yamlKeyString struct {
	yamlKeyExtendedCommon `yaml:",inline"`
	Value                 String
}

type yamlKeyData struct {
	yamlKeyExtendedCommon `yaml:",inline"`
	Value                 Data
}

type yamlKeyInteger struct {
	yamlKeyExtendedCommon `yaml:",inline"`
	Value                 Integer
}

type yamlKeyFloat struct {
	yamlKeyExtendedCommon `yaml:",inline"`
	Value                 Float
}

type yamlKeyBoolean struct {
	yamlKeyExtendedCommon `yaml:",inline"`
	Value                 Boolean
}
type yamlKeyDate struct {
	yamlKeyExtendedCommon `yaml:",inline"`
	Value                 Date
}
type yamlKeyArray struct {
	yamlKeyExtendedCommon `yaml:",inline"`
	Value                 Array
}
type yamlKeyArrayAdd struct {
	yamlKeyExtendedCommon `yaml:",inline"`
	Value                 ArrayAdd
}
type yamlKeyDict struct {
	yamlKeyExtendedCommon `yaml:",inline"`
	Value                 Dict
}
type yamlKeyDictAdd struct {
	yamlKeyExtendedCommon `yaml:",inline"`
	Value                 DictAdd
}

func (k *Key) UnmarshalYAML(value *yaml.Node) error {
	var auxInline yamlKeyInline
	if err := value.Decode(&auxInline); err == nil {
		k.Type = StringType
		nv := String(auxInline)
		k.Value = &nv
		return nil
	}
	var auxExtendedType yamlKeyExtendedCommon
	if err := value.Decode(&auxExtendedType); err != nil {
		return err
	}
	switch auxExtendedType.Type {
	case StringType:
		var auxExtended yamlKeyString
		if err := value.Decode(&auxExtended); err != nil {
			return err
		}
		k.Value = &auxExtended.Value
	case DataType:
		var auxExtended yamlKeyData
		if err := value.Decode(&auxExtended); err != nil {
			return err
		}
		k.Value = &auxExtended.Value
	case IntegerType:
		var auxExtended yamlKeyInteger
		if err := value.Decode(&auxExtended); err != nil {
			return err
		}
		k.Value = &auxExtended.Value
	case FloatType:
		var auxExtended yamlKeyFloat
		if err := value.Decode(&auxExtended); err != nil {
			return err
		}
		k.Value = &auxExtended.Value
	case BooleanType:
		var auxExtended yamlKeyBoolean
		if err := value.Decode(&auxExtended); err != nil {
			return err
		}
		k.Value = &auxExtended.Value
	case DateType:
		var auxExtended yamlKeyDate
		if err := value.Decode(&auxExtended); err != nil {
			return err
		}
		k.Value = &auxExtended.Value
	case ArrayType:
		var auxExtended yamlKeyArray
		if err := value.Decode(&auxExtended); err != nil {
			return err
		}
		k.Value = &auxExtended.Value
	case ArrayAddType:
		var auxExtended yamlKeyArrayAdd
		if err := value.Decode(&auxExtended); err != nil {
			return err
		}
		k.Value = &auxExtended.Value
	case DictType:
		var auxExtended yamlKeyDict
		if err := value.Decode(&auxExtended); err != nil {
			return err
		}
		k.Value = &auxExtended.Value
	case DictAddType:
		var auxExtended yamlKeyDictAdd
		if err := value.Decode(&auxExtended); err != nil {
			return err
		}
		k.Value = &auxExtended.Value
	default:
		return fmt.Errorf("invalid type '%s'", k.Type)
	}
	return nil
}