package rcv

import (
	"fmt"
	"strings"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/u"
)

func listenIpProto(ips []string) (listeners []Listener) {
	if len(ips) > 0 {
		for i := range ips {
			ip := ips[i]
			if !u.IsIPv4(ip) && !strings.Contains(ip, "[") {
				ip = fmt.Sprintf("[%v]", ip)
			}

			if cnf.C.Service.Listen.Tcp {
				ln, err := NewListenerTcp(ip)
				if err != nil {
					log.ErrorS(
						"listener",
						fmt.Sprintf("Failed to bind to tcp://%v:%v - %v", ip, cnf.C.Service.Listen.Port, err),
					)
				} else {
					listeners = append(listeners, ln)
					log.Info("listener", fmt.Sprintf("Bound to tcp://%v:%v", ip, cnf.C.Service.Listen.Port))
				}
			}

			if cnf.C.Service.Listen.Udp {
				ln, err := NewListenerUdp(ip)
				if err != nil {
					log.ErrorS(
						"listener",
						fmt.Sprintf("Failed to bind to udp://%v:%v - %v", ip, cnf.C.Service.Listen.Port, err),
					)
				} else {
					listeners = append(listeners, ln)
					log.Info("listener", fmt.Sprintf("Bound to udp://%v:%v", ip, cnf.C.Service.Listen.Port))
				}
			}

		}
	}
	return
}

func Start() (listeners []Listener) {
	protoListeners := listenIpProto(cnf.C.Service.Listen.IP4)
	listeners = append(listeners, protoListeners...)

	protoListeners = listenIpProto(cnf.C.Service.Listen.IP6)
	listeners = append(listeners, protoListeners...)
	return
}
