package rcv

import (
	"fmt"
	"strings"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
)

func listenIpProto(ips []string) []Listener {
	listeners := []Listener{}
	if len(ips) > 0 {
		for i := range ips {
			ip := ips[i]
			if strings.Contains(ip, ":") && !strings.Contains(ip, "[") {
				ip = fmt.Sprintf("[%v]", ip)
			}

			if cnf.ListenTcp {
				ln, err := NewListenerTcp(ip)
				if err != nil {
					log.ErrorS(
						"listener",
						fmt.Sprintf("Failed to bind to tcp://%v:%v - %v", ip, cnf.ListenPort, err),
					)
				} else {
					listeners = append(listeners, ln)
					log.Info("listener", fmt.Sprintf("Bound to tcp://%v:%v", ip, cnf.ListenPort))
				}
			}

			if cnf.ListenUdp {
				ln, err := NewListenerUdp(ip)
				if err != nil {
					log.ErrorS(
						"listener",
						fmt.Sprintf("Failed to bind to udp://%v:%v - %v", ip, cnf.ListenPort, err),
					)
				} else {
					listeners = append(listeners, ln)
					log.Info("listener", fmt.Sprintf("Bound to udp://%v:%v", ip, cnf.ListenPort))
				}
			}

		}
	}

	return listeners
}

func Start() []Listener {
	listeners := []Listener{}
	protoListeners := listenIpProto(cnf.ListenIP4)
	listeners = append(listeners, protoListeners...)

	protoListeners = listenIpProto(cnf.ListenIP6)
	listeners = append(listeners, protoListeners...)
	return listeners
}
