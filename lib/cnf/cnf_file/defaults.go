package cnf_file

import "github.com/superstes/calamary/cnf"

func applyConfigDefaults(newCnf *cnf.Config) {
	for i := range newCnf.Service.Listen {
		applyListenerDefaults(&newCnf.Service.Listen[i])
	}
	applyTimeoutDefaults(newCnf)
	applyMetricsDefaults(newCnf)
	applyOutputDefaults(newCnf)
}

func applyListenerDefaults(lncnf *cnf.ServiceListener) {
	if len(lncnf.IP4) == 0 && len(lncnf.IP6) == 0 {
		lncnf.IP4 = []string{"127.0.0.1"}
		lncnf.IP6 = []string{"::1"}
	}
}

func applyTimeoutDefaults(newCnf *cnf.Config) {
	if newCnf.Service.Timeout.Connect == 0 {
		newCnf.Service.Timeout.Connect = cnf.DefaultTimeoutConnect
	}
	if newCnf.Service.Timeout.Process == 0 {
		newCnf.Service.Timeout.Process = cnf.DefaultTimeoutProcess
	}
	if newCnf.Service.Timeout.Idle == 0 {
		newCnf.Service.Timeout.Idle = cnf.DefaultTimeoutIdle
	}
}

func applyMetricsDefaults(newCnf *cnf.Config) {
	if newCnf.Service.Metrics.Port == 0 {
		newCnf.Service.Metrics.Port = cnf.DefaultMetricsPort
	}
}

func applyOutputDefaults(newCnf *cnf.Config) {
	if newCnf.Service.Output.Retries == 0 {
		newCnf.Service.Output.Retries = cnf.DefaultConnectRetries
	}
}
