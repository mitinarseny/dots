package defaults

import "gopkg.in/yaml.v3"

const (
	AppType    = "app"
	DomainType = "domain"
	GlobalType = "global"
)

type Defaults struct {
	Apps    Domains
	Domains Domains
	Globals Domain
}

type yamlDefaults struct {
	Apps    Domains
	Domains Domains
	Globals Domain
}

func (d *Defaults) UnmarshalYAML(value yaml.Node) error {
	var aux yamlDefaults
	if err := value.Decode(&aux); err != nil {
		return err
	}

	d.Apps.Type = "apps"
	d.Apps = aux.Apps

	d.Domains.Type = "domains"

	return nil
}
