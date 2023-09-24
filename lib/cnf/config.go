package cnf

var DEBUG bool = false
var LOG_TIME bool = true
var C = Config{}
var RULES = []Rule{}

type Config struct {
	Service ServiceConfig `yaml:"service"`
	Vars    []Var         `yaml:"vars"`
	Rules   []RuleRaw     `yaml:"rules"`
}

type ServiceConfig struct {
	Timeout ServiceConfigTimeout `yaml:"timeout"`
	Listen  ServiceConfigListen  `yaml:"listen"`
	Output  ServiceConfigOutput  `yaml:"output"`
}

type ServiceConfigListen struct {
	Port        int      `yaml:"port" default="4128"`
	IP4         []string `ip4:"ip4"`
	IP6         []string `ip6:"ip6"`
	Tcp         bool     `yaml:"tcp" default="true"`
	Udp         bool     `yaml:"udp" default="false"` // not implemented
	Transparent bool     `yaml:"transparent" "default=false"`
}

type ServiceConfigTimeout struct {
	Connection int `yaml:"connection" default="5"`
	Handshake  int `yaml:"handshake" default="5"`
	Dial       int `yaml:"dial" default="5"`
	Intercept  int `yaml:"intercept" default="2"`
}

type ServiceConfigOutput struct {
	FwMark    int    `yaml:"fwmark" default="0"`
	Interface string `yaml:"interface" default=""`
}
