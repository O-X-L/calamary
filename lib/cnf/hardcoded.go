package cnf

import (
	"time"
)

const (
	VERSION         string = "1.0"
	LOG_TIME_FORMAT string = "2006-01-02 15:04:05"
	UDP_TTL                = 30 * time.Second
	UDP_BUFFER_SIZE int    = 4096
	BYTES_HDR_L5    int    = 5
	BYTES_TLS_REC   int    = 5
	BYTES_TLS_HS    int    = 4
	BYTES_TLS_EXT   int    = 4
)

var ConfigFileAbs string = "/etc/calamary/config.yml"
