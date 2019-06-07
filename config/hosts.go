package config

import "github.com/mitinarseny/dots/config/defaults"

type Host struct {
	Variables Variables
	Links Links
	Commands Commands
	Defaults defaults.Defaults
}

//type yamlHost struct {
//	Variables []struct {
//		Name  string
//		Value string
//	}
//	Links map[string]struct {
//		Path  string
//		Force bool
//	}
//	Commands []string
//	Defaults struct {
//		Apps, Domains, Globals map[string]map[string]interface{}
//	}
//}

//func (h *Host) UnmarshalYAML(value *yaml.Node) error {
//	var aux yamlHost
//	if err := value.Content.Decode(&aux); err != nil {
//		return err
//	}
//	varMapper := func(s string) string {
//		return h.Variables[s]
//	}
//	for _, v := range aux.Variables {
//		h.Variables[v.Name] = os.Expand(v.Value, varMapper)
//	}
//
//	ch := make(chan *Link)
//	for t, s := range aux.Links {
//		go func(target, source string, force bool) {
//			l := new(Link)
//			l.Target.Original = os.Expand(target, varMapper)
//			l.Source.Original = os.Expand(source, varMapper)
//			l.Force = force
//
//			ch <- l
//		}(t, s.Path, s.Force)
//	}
//
//	for range aux.Links {
//		h.Links = append(h.Links, <-ch)
//	}
//
//	for _, v := range aux.Commands {
//		cmd := Command(v)
//		h.Commands = append(h.Commands, &cmd)
//	}
//
//	apps, err := scanDefaults(aux.Defaults.Apps)
//	if err != nil {
//		return err
//	}
//	h.Defaults.Apps = apps
//
//	domains, err := scanDefaults(aux.Defaults.Domains)
//	if err != nil {
//		return err
//	}
//	h.Defaults.Domains = domains
//
//	globals, err := scanDefaults(aux.Defaults.Globals)
//	if err != nil {
//		return err
//	}
//	h.Defaults.Globals = globals
//	return nil
//}