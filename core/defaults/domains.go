package defaults

type Domain struct {
	// app, global,

	Keys []*Key
}

type yamlDomain map[string]Key

type Domains struct {
	Type string
	Domains map[string]Domain
}