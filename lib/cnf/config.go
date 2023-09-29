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
	Metrics MetricsConfig        `yaml:"metrics"`
}

type ServiceConfigListen struct {
	Port   int      `yaml:"port" default="4128"`
	IP4    []string `ip4:"ip4"`
	IP6    []string `ip6:"ip6"`
	Tcp    bool     `yaml:"tcp" default="true"`
	Udp    bool     `yaml:"udp" default="false"` // not implemented
	TProxy bool     `yaml:"tproxy" "default=false"`
}

// todo: defaults not working; set to 0
type ServiceConfigTimeout struct {
	Connect uint `yaml:"connect" default="2000"` // dial
	Process uint `yaml:"process" default="1000"` // parsing packet
	Idle    uint `yaml:"idle" default="15000"`   // close connection if no data was sent or received
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

type MetricsConfig struct {
	Enabled bool `yaml:"enabled" default="false"`
	Port    int  `yaml:"port" default="9512"`
}

// shortcut to setting as it is referenced often
func Metrics() bool {
	return C.Service.Metrics.Enabled
}
