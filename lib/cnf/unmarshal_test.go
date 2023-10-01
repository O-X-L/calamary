package cnf

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

func TestConfigDefaults(t *testing.T) {
	cnfListen1 := ServiceListener{}
	if err := defaults.Set(&cnfListen1); err != nil {
		t.Errorf("Unmarshal defaults - unable to set defaults #1! (%v)", err)
	}

	if len(cnfListen1.IP4) != 1 || cnfListen1.IP4[0] != "127.0.0.1" {
		t.Errorf("Unmarshal defaults-service #1 (%+v)", cnfListen1)
	}
	if len(cnfListen1.IP6) != 1 || cnfListen1.IP6[0] != "::1" {
		t.Errorf("Unmarshal defaults-service #2 (%+v)", cnfListen1)
	}
	if !cnfListen1.Tcp {
		t.Errorf("Unmarshal defaults-service #3 (%+v)", cnfListen1)
	}

	// test with acutal config file
	_, pathToTest, _, _ := runtime.Caller(0)
	pathToTestConfig := filepath.Dir(pathToTest) + "/testdata/unmarshal_test_defaults.yml"

	cnf2 := Config{}
	configRaw, err := os.ReadFile(pathToTestConfig)
	if err != nil {
		t.Errorf("Unmarshal defaults - unable to load test-config! (%v)", err)
	}
	err = yaml.Unmarshal(configRaw, &cnf2)
	if err != nil {
		t.Errorf("Unmarshal defaults - unable to unmarshal test-config! (%v)", err)
	}

	if len(cnf2.Service.Listen[0].IP4) != 1 || cnf2.Service.Listen[0].IP4[0] != "127.0.0.1" {
		t.Errorf("Unmarshal defaults-file #1 (%+v)", cnf2.Service.Listen[0])
	}
	if len(cnf2.Service.Listen[0].IP6) != 1 || cnf2.Service.Listen[0].IP6[0] != "::1" {
		t.Errorf("Unmarshal defaults-file #2 (%+v)", cnf2.Service.Listen[0])
	}
	if !cnf2.Service.Listen[0].Tcp {
		t.Errorf("Unmarshal defaults-file #3 (%+v)", cnf2.Service.Listen[0])
	}
	if cnf2.Service.Timeout.Connect != 2000 ||
		cnf2.Service.Timeout.Process != 1000 ||
		cnf2.Service.Timeout.Idle != 30000 {
		t.Errorf("Unmarshal defaults-file #4 (%+v)", cnf2.Service.Timeout)
	}
	if cnf2.Service.Output.Retries != 1 || cnf2.Service.Output.FwMark != 0 {
		t.Errorf("Unmarshal defaults-file #5 (%+v)", cnf2.Service.Output)
	}
	if cnf2.Service.Metrics.Port != 9512 {
		t.Errorf("Unmarshal defaults-file #6 (%+v)", cnf2.Service.Metrics)
	}
}
