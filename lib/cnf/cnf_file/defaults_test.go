package cnf_file

import (
	"testing"

	"github.com/superstes/calamary/cnf"
)

func TestApplyListenerDefaults(t *testing.T) {
	ln1 := cnf.ServiceListener{}
	if len(ln1.IP4) != 0 || len(ln1.IP6) != 0 {
		t.Error("Timeout config-defaults #1")
	}
	applyListenerDefaults(&ln1)
	if len(ln1.IP4) != 1 ||
		ln1.IP4[0] != "127.0.0.1" ||
		len(ln1.IP6) != 1 ||
		ln1.IP6[0] != "::1" {
		t.Error("Timeout config-defaults #2")
	}
}

func TestApplyTimeoutDefaults(t *testing.T) {
	cnf1 := cnf.Config{}
	if cnf1.Service.Timeout.Connect != 0 ||
		cnf1.Service.Timeout.Process != 0 ||
		cnf1.Service.Timeout.Idle != 0 {
		t.Error("Timeout config-defaults #1")
	}
	applyTimeoutDefaults(&cnf1)
	if cnf1.Service.Timeout.Connect != cnf.DefaultTimeoutConnect ||
		cnf1.Service.Timeout.Process != cnf.DefaultTimeoutProcess ||
		cnf1.Service.Timeout.Idle != cnf.DefaultTimeoutIdle {
		t.Error("Timeout config-defaults #2")
	}
}

func TestApplyMetricsDefaults(t *testing.T) {
	cnf1 := cnf.Config{}
	if cnf1.Service.Metrics.Port != 0 {
		t.Error("Timeout config-defaults #1")
	}
	applyMetricsDefaults(&cnf1)
	if cnf1.Service.Metrics.Port != cnf.DefaultMetricsPort {
		t.Error("Timeout config-defaults #2")
	}
}

func TestApplyOutputDefaults(t *testing.T) {
	cnf1 := cnf.Config{}
	if cnf1.Service.Output.Retries != 0 {
		t.Error("Timeout config-defaults #1")
	}
	applyOutputDefaults(&cnf1)
	if uint16(cnf1.Service.Output.Retries) != uint16(cnf.DefaultConnectRetries) {
		t.Error("Timeout config-defaults #2")
	}
}
