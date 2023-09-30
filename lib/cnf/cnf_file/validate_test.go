package cnf_file

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/superstes/calamary/cnf"
)

func TestValidateListener(t *testing.T) {
	if !validateListener(cnf.ServiceListener{Mode: "transparent", Port: 4128}, false) {
		t.Error("Listener config-validation #1")
	}
	if !validateListener(cnf.ServiceListener{Mode: "proxyproto", Port: 4128}, false) {
		t.Error("Listener config-validation #2")
	}
	if !validateListener(cnf.ServiceListener{Mode: "http", Port: 4128}, false) {
		t.Error("Listener config-validation #3")
	}
	if !validateListener(cnf.ServiceListener{Mode: "https", Port: 4128}, false) {
		t.Error("Listener config-validation #4")
	}
	if !validateListener(cnf.ServiceListener{Mode: "socks5", Port: 4128}, false) {
		t.Error("Listener config-validation #5")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Listener config-validation #6")
		}
	}()
	validateListener(cnf.ServiceListener{Mode: "random", Port: 4128}, true)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Listener config-validation #7")
		}
	}()
	validateListener(cnf.ServiceListener{Mode: "transparent"}, true)
}

func TestValidateCerts(t *testing.T) {
	_, pathToTest, _, _ := runtime.Caller(0)
	pathToTestCerts := filepath.Dir(pathToTest) + "/testdata/"

	if !validateCerts(cnf.ServiceCertificates{}, false) {
		t.Error("Certificate config-validation #1")
	}

	if !validateCerts(cnf.ServiceCertificates{
		ServerPublic:  pathToTestCerts + "server.crt",
		ServerPrivate: pathToTestCerts + "server.key",
	}, false) {
		t.Error("Certificate config-validation #2")
	}

	if validateCerts(cnf.ServiceCertificates{
		ServerPublic: pathToTestCerts + "server.crt",
	}, false) {
		t.Errorf("Certificate config-validation #3")
	}

	if validateCerts(cnf.ServiceCertificates{
		ServerPrivate: pathToTestCerts + "server.key",
	}, false) {
		t.Errorf("Certificate config-validation #4")
	}

	if validateCerts(cnf.ServiceCertificates{
		ServerPrivate: pathToTestCerts + "server.crt",
		ServerPublic:  pathToTestCerts + "server.crt",
	}, false) {
		t.Errorf("Certificate config-validation #5")
	}

	if validateCerts(cnf.ServiceCertificates{
		ServerPrivate: pathToTestCerts + "server.key",
		ServerPublic:  pathToTestCerts + "server.key",
	}, false) {
		t.Errorf("Certificate config-validation #6")
	}
}
