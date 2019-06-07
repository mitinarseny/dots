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
	//
	//if err := value.Decode(&k.Value); err != nil {
	//	return err
	//}
	//k.Value = auxExtendedType.Value.(interface{
	//	fmt.Stringer
	//})

	return nil

	//var auxExtendedType yamlKeyExtended
	//if err := value.Decode(&auxExtendedType); err == nil {
	//	k.Type = auxExtendedType.Type
	//	switch k.Type {
	//	case StringType:
	//		k.Value = new(String)
	//	case DataType:
	//		k.Value = new(Data)
	//	case IntegerType:
	//		k.Value = new(Integer)
	//	case FloatType:
	//		k.Value = new(Float)
	//	case BooleanType:
	//		k.Value = new(Boolean)
	//	case DateType:
	//		k.Value = new(Date)
	//	case ArrayType:
	//		k.Value = new(Array)
	//	case ArrayAddType:
	//		k.Value = new(ArrayAdd)
	//	case DictType:
	//		k.Value = new(Dict)
	//	case DictAddType:
	//		k.Value = new(DictAdd)
	//	}
	//
	//	k.Value = auxExtendedType.Value
	//	return nil
	//}
	//return nil
}

//type Scanner interface {
//	Scan(interface{}) error
//}
//
////type Key struct {
////	Type  string
////	Value interface {
////		fmt.Stringer
////		Scanner
////	}
////}
//
//type scanKey struct {
//	Type  string
//	Value interface{}
//}
//
//func (r *Key) String() string {
//	return r.Value.String()
//}
//
//func (r *Key) Scan(i interface{}) error {
//	if r == nil || r.Type == "" {
//		switch v := i.(type) {
//		case string:
//			*r = Key{
//				Type:  StringType,
//				Value: new(String),
//			}
//			if err := r.Value.Scan(i); err != nil {
//				return err
//			}
//		case scanKey:
//			*r = Key{
//				Type: v.Type,
//			}
//			if err := r.scanWithType(&v); err != nil {
//				return nil
//			}
//		default:
//			return fmt.Errorf("cannot scan %T from %T", r, v)
//		}
//	}
//	return nil
//
//}
//
//func (r *Key) scanWithType(s *scanKey) error {
//	switch r.Type {
//	case StringType:
//		r.Value = new(String)
//	case DataType:
//		r.Value = new(Data)
//	case IntegerType:
//		r.Value = new(Integer)
//	case FloatType:
//		r.Value = new(Float)
//	case BooleanType:
//		r.Value = new(Boolean)
//	case DateType:
//		r.Value = new(Date)
//	case ArrayType:
//		r.Value = new(Array)
//	case ArrayAddType:
//		r.Value = new(ArrayAdd)
//	case DictType:
//		r.Value = new(Dict)
//	case DictAddType:
//		r.Value = new(DictAdd)
//	default:
//		return fmt.Errorf("no type '%s'", r.Type)
//	}
//	if err := r.Value.Scan(s.Value); err != nil {
//		return err
//	}
//	return nil
//}
