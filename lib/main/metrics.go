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

func catchAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

func startPrometheusExporter() {
	if cnf.Metrics() {
		log.Info("service", "Starting prometheus metrics-exporter")
		prometheus.Register(metrics.BytesRcv)
		prometheus.Register(metrics.BytesSent)
		prometheus.Register(metrics.CurrentConn)
		prometheus.Register(metrics.ReqTcp)
		prometheus.Register(metrics.RuleReqL3Proto)
		prometheus.Register(metrics.RuleReqL5Proto)
		prometheus.Register(metrics.RuleReqTlsVersion)
		prometheus.Register(metrics.RuleHits)
		prometheus.Register(metrics.RuleMatches)
		prometheus.Register(metrics.RuleActions)

		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/", catchAll)
		http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", cnf.C.Service.Metrics.Port), nil)
		http.ListenAndServe(fmt.Sprintf("[::1]:%v", cnf.C.Service.Metrics.Port), nil)

		prometheus.MustRegister(metrics.BytesRcv)
		prometheus.MustRegister(metrics.BytesSent)
		prometheus.MustRegister(metrics.CurrentConn)
		prometheus.MustRegister(metrics.ReqTcp)
		prometheus.MustRegister(metrics.RuleReqL3Proto)
		prometheus.MustRegister(metrics.RuleReqL5Proto)
		prometheus.MustRegister(metrics.RuleReqTlsVersion)
		prometheus.MustRegister(metrics.RuleHits)
		prometheus.MustRegister(metrics.RuleMatches)
		prometheus.MustRegister(metrics.RuleActions)
	}
}
