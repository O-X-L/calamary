package cnf

import (
	"github.com/superstes/calamary/proc/meta"
)

var LOG_TIME bool = true
var C *Config
var RULES *[]Rule

type Config struct {
	Service ServiceConfig `yaml:"service"`
	Vars    []Var         `yaml:"vars"`
	Rules   []RuleRaw     `yaml:"rules"`
}

type ServiceConfig struct {
	Timeout ServiceTimeout      `yaml:"timeout"`
	Listen  []ServiceListener   `yaml:"listen"`
	Certs   ServiceCertificates `yaml:"certs"`
	Output  ServiceOutput       `yaml:"output"`
	Debug   bool                `yaml:"debug" default:"false"`
	Metrics ServiceMetrics      `yaml:"metrics"`
}

// todo: implement default listen-ips = localhost
// todo: make sure mode is valid
// todo: if no listeners were provided - start only transparent
type ServiceListener struct {
	Mode   meta.ListenMode `yaml:"mode" default:"transparent"`
	Port   uint16          `yaml:"port"`
	IP4    []string        `yaml:"ip4" default:"[\"127.0.0.1\"]"`
	IP6    []string        `yaml:"ip6" default:"[\"::1\"]"`
	Tcp    bool            `yaml:"tcp" default:"true"`
	Udp    bool            `yaml:"udp" default:"false"` // not implemented
	TProxy bool            `yaml:"tproxy" default:"false"`
}

type ServiceTimeout struct {
	Connect uint `yaml:"connect" default:"2000"` // dial
	Process uint `yaml:"process" default:"1000"` // parsing packet
	Idle    uint `yaml:"idle" default:"30000"`   // close connection if no data was sent or received
}

var DefaultConnectRetryWait = uint(1000) // ms

type ServiceOutput struct {
	FwMark    uint8  `yaml:"fwmark" default:"0"`
	Interface string `yaml:"interface" default:""`
	// IP4       []string `ip4:"ip4"`
	// IP6       []string `ip6:"ip6"`
	Retries uint8 `yaml:"retries" default:"1"`
}

var DefaultMetricsPort = uint16(9512)

type ServiceMetrics struct {
	Enabled bool   `yaml:"enabled" default:"false"`
	Port    uint16 `yaml:"port" default:"9512"`
}

type ServiceCertificates struct {
	ServerPublic     string `yaml:"serverPublic"`
	ServerPrivate    string `yaml:"serverPrivate"`
	InterceptPublic  string `yaml:"interceptPublic"`
	InterceptPrivate string `yaml:"interceptPrivate"`
}

// shortcut to setting as it is referenced often
func Metrics() bool {
	return C.Service.Metrics.Enabled
}

// shortcut to setting as it is referenced often
func Debug() bool {
	return C.Service.Debug
}
