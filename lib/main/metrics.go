package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/metrics"
)

// todo: clean-up redundant linking
var metricFuncs = []prometheus.Collector{
	metrics.BytesRcv,
	metrics.BytesSent,
	metrics.CurrentConn,
	metrics.ReqTcp,
	metrics.ReqL3Proto,
	metrics.ReqL5Proto,
	metrics.ReqTlsVersion,
	metrics.RuleHits,
	metrics.RuleMatches,
	metrics.RuleActions,
}

func denyAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

func startPrometheusExporter() {
	if cnf.Metrics() {
		log.Info("service", "Starting prometheus metrics-exporter")

		for _, mf := range metricFuncs {
			prometheus.Register(mf)
		}

		metricsSrv := http.NewServeMux()
		metricsSrv.Handle("/metrics", promhttp.Handler())
		metricsSrv.HandleFunc("/", denyAll)
		http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", cnf.C.Service.Metrics.Port), metricsSrv)
		http.ListenAndServe(fmt.Sprintf("[::1]:%v", cnf.C.Service.Metrics.Port), metricsSrv)

		for _, mf := range metricFuncs {
			prometheus.MustRegister(mf)
		}
	}
}
