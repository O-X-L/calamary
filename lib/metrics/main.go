package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	BytesRcv = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calamary",
		Name:      "bytes_rcv",
		Help:      "The total number of bytes received",
	})
	BytesSent = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calamary",
		Name:      "bytes_sent",
		Help:      "The total number of bytes sent",
	})
	CurrentConn = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "calamary",
		Name:      "current_connections",
		Help:      "Number of currently active connections",
	})
	ReqTcp = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "calamary",
		Name:      "req_tcp",
		Help:      "Count of tcp requests received",
	})
	RuleReqL3Proto = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "calamary",
			Name:      "req_protoL3",
			Help:      "Count of requests per L3 protocol",
		},
		[]string{"protocol"},
	)
	RuleReqL5Proto = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "calamary",
			Name:      "req_protoL5",
			Help:      "Count of requests per L5 protocol",
		},
		[]string{"protocol"},
	)
	RuleReqTlsVersion = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "calamary",
			Name:      "req_tls_version",
			Help:      "Count of requests per TLS version",
		},
		[]string{"version"},
	)
	RuleHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "calamary",
			Name:      "rule_hits",
			Help:      "Times a rule got evaluated",
		},
		[]string{"ruleId"},
	)
	RuleMatches = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "calamary",
			Name:      "rule_matches",
			Help:      "Times traffic matched a rule",
		},
		[]string{"ruleId"},
	)
	RuleActions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "calamary",
			Name:      "rule_actions",
			Help:      "Times an rule-action got applied",
		},
		[]string{"action"},
	)
)
