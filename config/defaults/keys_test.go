package defaults

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestKeyInline(t *testing.T) {
	r := require.New(t)

	data := `value`
	expectedValue := String("value")

	k := new(Key)
	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, StringType)
	r.IsType(new(String), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}

func TestKeyString(t *testing.T) {
	r := require.New(t)

	data := `type: string
value: value`
	expectedValue := String("value")

	k := new(Key)
	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, StringType)
	r.IsType(new(String), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}

func TestKeyData(t *testing.T) {
	r := require.New(t)

	data := `type: data
value: value`
	expectedValue := Data("value")

	k := new(Key)
	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, DataType)
	r.IsType(new(Data), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}

func TestKeyInteger(t *testing.T) {
	r := require.New(t)

	data := `type: int
value: 1`
	expectedValue := Integer(1)
	k := new(Key)

	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, IntegerType)
	r.IsType(new(Integer), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}

func TestKeyFloat(t *testing.T) {
	r := require.New(t)

	data := `type: float
value: 1.1`
	expectedValue := Float(1.1)
	k := new(Key)

	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, FloatType)
	r.IsType(new(Float), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}

func TestKeyBoolean(t *testing.T) {
	r := require.New(t)

	data := `type: bool
value: true`
	expectedValue := Boolean(true)
	k := new(Key)

	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, BooleanType)
	r.IsType(new(Boolean), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}

func TestKeyDate(t *testing.T) {
	r := require.New(t)

	data := `type: date
value: 1999-12-20 04:00:00 +0000`
	expectedValue := Date("1999-12-20 04:00:00 +0000")
	k := new(Key)

	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, DateType)
	r.IsType(new(Date), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}

func TestKeyArray(t *testing.T) {
	r := require.New(t)

	data := `type: array
value:
    - value1
    - value2
    - value3`
	expectedValue := Array{
		"value1",
		"value2",
		"value3",
	}
	k := new(Key)

	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, ArrayType)
	r.IsType(new(Array), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}

func TestKeyArrayAdd(t *testing.T) {
	r := require.New(t)

	data := `type: array-add
value:
    - value1
    - value2
    - value3`
	expectedValue := ArrayAdd{
		"value1",
		"value2",
		"value3",
	}
	k := new(Key)

	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, ArrayAddType)
	r.IsType(new(ArrayAdd), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}

func TestKeyDict(t *testing.T) {
	r := require.New(t)

	data := `type: dict
value:
    key1: value1
    key2: value2
    key3: value3`
	expectedValue := Dict{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	k := new(Key)

	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, DictType)
	r.IsType(new(Dict), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}

func TestKeyDictAdd(t *testing.T) {
	r := require.New(t)

	data := `type: dict-add
value:
    key1: value1
    key2: value2
    key3: value3`
	expectedValue := DictAdd{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	k := new(Key)

	err := yaml.Unmarshal([]byte(data), &k)

	r.NoError(err)
	r.Equal(k.Type, DictAddType)
	r.IsType(new(DictAdd), k.Value)
	r.EqualValues(&expectedValue, k.Value)
}
