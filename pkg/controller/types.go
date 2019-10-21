package controller

// Team ...
type Team struct {
	Team       string   `yaml:"team"`
	Value      string   `yaml:"value"`
	NameSpaces []string `yaml:"namespaces"`
}

// Config ...
type Config struct {
	Maintainers []Team `yaml:"maintainers"`
}
