package controller

// Team ...
type Team struct {
	Team       string            `yaml:"team"`
	Labels     map[string]string `yaml:"labels"`
	NameSpaces []string          `yaml:"namespaces"`
}

// Config ...
type Config struct {
	Maintainers []Team `yaml:"maintainers"`
}
