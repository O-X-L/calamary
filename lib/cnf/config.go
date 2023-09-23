package cnf

import (
	"time"
)

const (
	VERSION         string = "1.0"
	LOG_TIME_FORMAT        = "2006-01-02 15:04:05"
	UDP_TTL                = 30 * time.Second
	UDP_BUFFER_SIZE        = 4096
)

var DEBUG bool = false
var LOG_TIME bool = true
var CONFIG = Config{}

type Config struct {
	Timeout ConfigTimeout `yaml:"timeout"`
	Listen  ConfigListen  `yaml:"listen"`
	Output  ConfigOutput  `yaml:"output"`
}

type ConfigListen struct {
	Port        int      `yaml:"port,default=4128"`
	IP4         []string `ip4:"ip4"`
	IP6         []string `ip6:"ip6"`
	Tcp         bool     `yaml:"tcp,default=true"`
	Udp         bool     `yaml:"udp,default=false"`
	Transparent bool     `yaml:"transparent,default=false"`
}

type ConfigTimeout struct {
	Connection int `yaml:"connection,default=5"`
	Handshake  int `yaml:"handshake,default=5"`
	Dial       int `yaml:"dial,default=5"`
	Intercept  int `yaml:"intercept,default=5"`
}

type ConfigOutput struct {
	FwMark    int `yaml:"fwmark,default=0"`
	Interface int `yaml:"interface"`
}
