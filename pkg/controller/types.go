package controller

// Team ...
type Team struct {
	Team       string            `yaml:"team"`
	Labels     map[string]string `yaml:"labels"`
	NameSpaces []string          `yaml:"namespaces"`
}

// Config ...
type Config struct {
	Maintainers     []Team   `yaml:"maintainers"`
	AdminNamespaces []string `yaml:"adminnamespaces"`
	LimitCPU        string   `yaml:"limitcpu"`
	LimitMemory     string   `yaml:"limitmemory"`
	RequestCPU      string   `yaml:"requestcpu"`
	RequestMemory   string   `yaml:"requestmemory"`
}
