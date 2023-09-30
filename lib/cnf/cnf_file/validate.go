package cnf_file

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
)

func validateConfig(newCnf cnf.Config, fail bool) bool {
	for i := range newCnf.Service.Listen {
		if !validateListener(newCnf.Service.Listen[i], fail) {
			return false
		}
	}
	if !validateCerts(newCnf.Service.Certs, fail) {
		return false
	}
	return true
}

func validateListener(lncnf cnf.ServiceListener, fail bool) bool {
	if lncnf.Mode != meta.ListenModeTransparent && lncnf.Mode != meta.ListenModeProxyProto &&
		lncnf.Mode != meta.ListenModeHttp && lncnf.Mode != meta.ListenModeHttps &&
		lncnf.Mode != meta.ListenModeSocks5 {
		if !fail {
			return false
		}
		panic(fmt.Errorf("Listener mode '%s' is not valid", lncnf.Mode))
	}
	if lncnf.Port == 0 {
		if !fail {
			return false
		}
		panic(fmt.Errorf("No listen-port supplied for mode '%s'", lncnf.Mode))
	}
	return true
}

func validateCerts(certs cnf.ServiceCertificates, fail bool) bool {
	if certs.ServerPublic != "" || certs.ServerPrivate != "" {
		if !validateCert(certs.ServerPublic, true, fail) {
			return false
		}
		if !validateCert(certs.ServerPrivate, false, fail) {
			return false
		}
	}
	// todo: check if valid sub-ca
	if certs.InterceptPublic != "" || certs.InterceptPrivate != "" {
		if !validateCert(certs.InterceptPublic, true, fail) {
			return false
		}
		if !validateCert(certs.InterceptPrivate, false, fail) {
			return false
		}
	}
	return true
}

func validateCert(file string, public bool, fail bool) bool {
	raw, err := os.ReadFile(file)
	if err != nil {
		if !fail {
			return false
		}
		panic(fmt.Errorf("Certificate file could not be loaded: '%s' (%v)", file, err))
	}
	var cert tls.Certificate
	for {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, block.Bytes)
		} else {
			cert.PrivateKey, err = parsePrivateKey(block.Bytes)
			if err != nil {
				if !fail {
					return false
				}
				panic(fmt.Errorf("Certificate private-key could not be parsed: '%s' (%v)", file, err))
			}
		}
		raw = rest
	}
	if (public && len(cert.Certificate) == 0) || (!public && cert.PrivateKey == nil) {
		if !fail {
			return false
		}
		panic(fmt.Errorf("Certificate file could not be loaded: '%s'", file))
	}
	return true
}

func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("Found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("Failed to parse private key")
}
