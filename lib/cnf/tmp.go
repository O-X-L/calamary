package cnf

import "time"

// config that will be move to the user-controlled config-file

var (
	ListenPort        int    = 4128
	ListenIP4                = []string{"127.0.0.1"}
	ListenIP6                = []string{"::1"}
	ListenTcp         bool   = true
	ListenUdp         bool   = false // not implemented
	ListenTransparent bool   = false
	TimeoutConnection        = 5 * time.Second
	TimeoutHandshake         = 5 * time.Second
	TimeoutDial              = 5 * time.Second
	TimeoutIntercept         = 2 * time.Second
	OutputFwMark      int    = 0
	OutputInterface   string = ""
)
