package cnf

// vars should parsed first so we can replace all their references when parsing other parts of the config

type Var struct {
	Name  string   `yaml:"name"`
	Value []string `yaml:"value"`
}
