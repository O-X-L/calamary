package cnf

import (
	"time"
)

const (
	VERSION         string = "1.0"
	LOG_TIME_FORMAT string = "2006-01-02 15:04:05"
	UDP_TTL                = 30 * time.Second
	UDP_BUFFER_SIZE int    = 4096
	CONFIG_FILE_ABS string = "/etc/calamary/config.yml"
)
