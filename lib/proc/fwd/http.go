package fwd

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/metrics"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
	"github.com/superstes/calamary/send"
	"github.com/superstes/calamary/u"
)

func ForwardHttp(srvCnf cnf.ServiceListener, l4Proto meta.Proto, conn net.Conn) {
	defer conn.Close()

	if l4Proto != meta.ProtoL4Tcp {
		log.ErrorS("forward", fmt.Sprintf("L4Proto %s not supported in http/https mode!", meta.RevProto(l4Proto)))
		return
	}

	if cnf.Metrics() {
		metrics.ReqTcp.Inc()
		metrics.CurrentConn.Inc()
		defer metrics.CurrentConn.Dec()
	}

	// todo: full parsing may not be needed if 'connect'
	pktProxy, connProxyIo := parseConn(srvCnf, l4Proto, conn)

	if pktProxy.L5.Proto != meta.ProtoL5Http {
		parse.LogConnError("forward", pktProxy, "Got non HTTP-Request on HTTP server")
		return
	}

	reqProxy := readRequest(pktProxy, connProxyIo)

	if reqProxy == nil {
		return
	} else if reqProxy.Method == http.MethodConnect {
		forwardConnect(srvCnf, l4Proto, conn, pktProxy, connProxyIo, reqProxy)
	} else {
		// plaintext HTTP is unsafe by design and should not be allowed
		// todo: implement generic https-redirection response

		forwardPlain(srvCnf, l4Proto, conn, pktProxy, connProxyIo, reqProxy)
	}
}

func forwardPlain(
	srvCnf cnf.ServiceListener, l4Proto meta.Proto, conn net.Conn,
	pkt parse.ParsedPacket, connIo io.ReadWriter, req *http.Request,
) {
	pkt.L5Http = parse.ParseHttpPacket(pkt, req)
	pkt.L4.DestPort = pkt.L5Http.Port
	dest := resolveTargetHostname(pkt)
	if dest == nil {
		proxyResp := responseFailed()
		proxyResp.Write(conn)
		return
	}
	pkt.L3.DestIP = dest
	parse.LogConnDebug("forward", pkt, "Updated destination IP")

	filterConn(pkt, conn, connIo)
	send.ForwardHttp(pkt, conn, connIo, req)
}

func forwardConnect(
	srvCnf cnf.ServiceListener, l4Proto meta.Proto, conn net.Conn,
	pktProxy parse.ParsedPacket, connProxyIo io.ReadWriter, reqProxy *http.Request,
) {
	proxyResp := http.Response{
		StatusCode: 200,
		Request:    reqProxy,
		ProtoMajor: reqProxy.ProtoMajor,
		ProtoMinor: reqProxy.ProtoMinor,
	}
	err := proxyResp.Write(conn)
	if err != nil {
		parse.LogConnError("forward", pktProxy, "Failed to write proxy response")
		return
	}

	pkt, connIo := parseConn(srvCnf, l4Proto, conn)
	host, port := parse.SplitHttpHost(reqProxy.Host, pkt.L5.Encrypted)
	pkt.L5Http = &parse.ParsedHttp{
		Host: host,
		Port: port,
	}
	pkt.L4.DestPort = pkt.L5Http.Port

	/*
		req := readRequest(pkt, connIo)
		pkt.L5Http = parse.ParseHttpPacket(pkt, req)
	*/
	dest := resolveTargetHostname(pkt)
	if dest == nil {
		proxyResp = responseFailed()
		proxyResp.Write(conn)
		return
	}
	pkt.L3.DestIP = dest
	parse.LogConnDebug("forward", pkt, "Updated destination IP")

	filterConn(pkt, conn, connIo)
	send.Forward(pkt, conn, connIo)

}

func resolveTargetHostname(pkt parse.ParsedPacket) net.IP {
	var destIp4, destIp6 []net.IP
	if pkt.L5.Encrypted == meta.OptBoolTrue {
		// todo: enable tls-interception
		if pkt.L5.TlsSni == "" {
			parse.LogConnError("forward", pkt, "Target hostname not retrievable via TLS-SNI")
			return nil
		}
		destIp4, destIp6 = u.DnsLookup46(pkt.L5.TlsSni)

	} else {
		if pkt.L5Http.Host == "" {
			parse.LogConnError("forward", pkt, "Target hostname not retrievable via Host-Header")
			return nil
		}
		destIp4, destIp6 = u.DnsLookup46(pkt.L5Http.Host)
	}

	dest := FirstReachableTarget(pkt, destIp4, destIp6, "tcp", pkt.L4.DestPort)
	if dest == nil {
		parse.LogConnError("send", pkt,
			fmt.Sprintf("No target IP reachable after %v retries: IP4 %v, IP6 %v, Port %v",
				cnf.C.Service.Output.Retries, destIp4, destIp6, pkt.L4.DestPort),
		)
		return nil
	}
	return dest
}

func responseFailed() http.Response {
	return http.Response{
		StatusCode: http.StatusRequestTimeout,
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
}

func readRequest(pkt parse.ParsedPacket, connIo io.ReadWriter) *http.Request {
	reqRaw := bufio.NewReader(connIo)
	req, err := http.ReadRequest(reqRaw)
	if err != nil {
		parse.LogConnError("forward", pkt, "Failed to parse HTTP request")
		return nil
	}
	return req
}
