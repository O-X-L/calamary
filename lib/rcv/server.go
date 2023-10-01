package rcv

import (
	"fmt"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/u"
)

var (
	serverModeMapping = map[meta.ListenMode]func(lncnf cnf.ServiceListener) []Server{
		meta.ListenModeTransparent: newServerTransparent,
		meta.ListenModeProxyProto:  newServerProxyProto,
		meta.ListenModeHttp:        newServerHttp,
		meta.ListenModeHttps:       newServerHttps,
		meta.ListenModeSocks5:      newServerSocks5,
	}
)

func newServersForIps(
	ips []string, lncnf cnf.ServiceListener,
	lnfuncTcp func(string, cnf.ServiceListener) (Server, error),
	lnfuncUdp func(string, cnf.ServiceListener) (Server, error),
) (servers []Server) {
	if len(ips) > 0 {
		for _, ip := range ips {
			ip = u.FormatIPv6(ip)
			if lncnf.Tcp && lnfuncTcp != nil {
				srv, err := lnfuncTcp(ip, lncnf)
				if err != nil {
					log.ErrorS("rcv", fmt.Sprintf("Failed to bind to tcp://%v:%v in %s mode - %v", ip, lncnf.Port, lncnf.Mode, err))
				} else {
					servers = append(servers, srv)
					log.Info("rcv", fmt.Sprintf("Bound to tcp://%v:%v in %s mode", ip, lncnf.Port, lncnf.Mode))
				}
			}
			if lncnf.Udp && lnfuncUdp != nil {
				srv, err := lnfuncUdp(ip, lncnf)
				if err != nil {
					log.ErrorS("rcv", fmt.Sprintf("Failed to bind to udp://%v:%v in %s mode - %v", ip, lncnf.Port, lncnf.Mode, err))
				} else {
					servers = append(servers, srv)
					log.Info("rcv", fmt.Sprintf("Bound to udp://%v:%v in %s mode", ip, lncnf.Port, lncnf.Mode))
				}
			}
		}
	}
	return
}

func newServerTransparent(lncnf cnf.ServiceListener) (servers []Server) {
	servers = newServersForIps(lncnf.IP4, lncnf, newServerTransparentTcp, serverNotImplemented)
	servers = append(servers, newServersForIps(lncnf.IP6, lncnf, newServerTransparentTcp, serverNotImplemented)...)
	return
}

func newServerProxyProto(lncnf cnf.ServiceListener) (servers []Server) {
	servers = newServersForIps(lncnf.IP4, lncnf, newServerProxyProtoTcp, serverNotImplemented)
	servers = append(servers, newServersForIps(lncnf.IP6, lncnf, newServerProxyProtoTcp, serverNotImplemented)...)
	return
}

func newServerHttp(lncnf cnf.ServiceListener) (servers []Server) {
	servers = newServersForIps(lncnf.IP4, lncnf, newServerHttpTcp, serverNotSupported)
	servers = append(servers, newServersForIps(lncnf.IP6, lncnf, newServerHttpTcp, serverNotSupported)...)
	return
}

func newServerHttps(lncnf cnf.ServiceListener) (servers []Server) {
	servers = newServersForIps(lncnf.IP4, lncnf, newServerHttpsTcp, serverNotSupported)
	servers = append(servers, newServersForIps(lncnf.IP6, lncnf, newServerHttpsTcp, serverNotSupported)...)
	return
}

func newServerSocks5(lncnf cnf.ServiceListener) (servers []Server) {
	servers = newServersForIps(lncnf.IP4, lncnf, serverNotImplemented, serverNotImplemented)
	servers = append(servers, newServersForIps(lncnf.IP6, lncnf, serverNotImplemented, serverNotImplemented)...)
	return
}

func BuildServers() (servers []Server) {
	for _, lncnf := range cnf.C.Service.Listen {
		servers = append(
			servers,
			serverModeMapping[lncnf.Mode](lncnf)...,
		)
	}
	log.Debug("rcv", fmt.Sprintf("SERVER DUMP: %+v", servers))
	return
}
