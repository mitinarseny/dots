package defaults

type Defaults struct {
	Apps Domains
	Domains Domains
	Globals Domains
}

type yamlDefaults map[string]map[string]interface{}

//func (d *Defaults) UnmarshalYAML(value *yaml.Node) error {
//
//}
