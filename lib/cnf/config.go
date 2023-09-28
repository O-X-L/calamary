package cnf

import "gopkg.in/yaml.v3"

var LOG_TIME bool = true
var C *Config
var RULES *[]Rule

type Config struct {
	Service ServiceConfig `yaml:"service"`
	Vars    []Var         `yaml:"vars"`
	Rules   []RuleRaw     `yaml:"rules"`
}

type ServiceConfig struct {
	Timeout ServiceConfigTimeout `yaml:"timeout"`
	Listen  ServiceConfigListen  `yaml:"listen"`
	Output  ServiceConfigOutput  `yaml:"output"`
	Debug   bool                 `yaml:"debug" default="false"`
}

type ServiceConfigListen struct {
	Port   int      `yaml:"port" default="4128"`
	IP4    []string `ip4:"ip4"`
	IP6    []string `ip6:"ip6"`
	Tcp    bool     `yaml:"tcp" default="true"`
	Udp    bool     `yaml:"udp" default="false"` // not implemented
	TProxy bool     `yaml:"tproxy" "default=false"`
}

type ServiceConfigTimeout struct {
	Connect uint `yaml:"connect" default="2000"`
	Process uint `yaml:"process" default="1000"`
}

type ServiceConfigOutput struct {
	FwMark    uint8  `yaml:"fwmark" default="0"`
	Interface string `yaml:"interface" default=""`
	// IP4       []string `ip4:"ip4"`
	// IP6       []string `ip6:"ip6"`
	// Retries   uint8  `yaml:"retries" default="1"`
}

// allow single string to be supplied
type YamlStringArray []string

func (a *YamlStringArray) UnmarshalYAML(value *yaml.Node) error {
	var multi []string
	err := value.Decode(&multi)
	if err != nil {
		var single string
		err := value.Decode(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}
