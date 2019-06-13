package core

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestLinks(t *testing.T) {
	r := require.New(t)

	data := `
~/.target1: .source1
~/.target2:
    path: .source2
    force: true`
	expected := Links{
		&Link{
			Source: FilePath{
				Original: ".source1",
			},
			Target: FilePath{
				Original: "~/.target1",
			},
			Force: false,
		},
		&Link{
			Source: FilePath{
				Original: ".source2",
			},
			Target: FilePath{
				Original: "~/.target2",
			},
			Force: true,
		},
	}
	l := new(Links)

	err := yaml.Unmarshal([]byte(data), l)

	r.NoError(err)
	r.EqualValues(&expected, l)
}