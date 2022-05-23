package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	_ "embed"
	"io/ioutil"
	"testing"

	"github.com/hkloudou/xlib/xcert"
)

//go:embed cert/ca.pem
var pem []byte

//go:embed cert/ca.key
var key []byte

func Test_domain(t *testing.T) {
	ca, caKey, err := xcert.ParseCertPair(pem, key)
	pem, key, err = xcert.GenerateEcdsaCertWithParent(
		xcert.Template(
			xcert.PkixName(pkix.Name{CommonName: "server"}),
			xcert.Hosts("localhost", "127.0.0.1"),
			xcert.IsCa(false),
			xcert.ExtKeyUsage(x509.ExtKeyUsageServerAuth),
		), ca, caKey)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("./cert/server.pem", pem, 0644)
	ioutil.WriteFile("./cert/server.key", key, 0644)
}
