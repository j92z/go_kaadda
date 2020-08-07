package main

import (
	"github/j92z/go_kaadda/genx"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	cert := genx.New(genx.CertConfig{
		Path: "certificate",
		Name: "kaadda",
		Host: []string{"127.0.0.1", "localhost"},
	})
	pool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(cert.RootCert.Path.Ca)
	if err != nil {
		log.Fatal("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(ca)
	fmt.Println(cert.ClientCert.Path.Ca, cert.ClientCert.Path.Key)
	cliCrt, err := tls.LoadX509KeyPair(cert.ClientCert.Path.Ca, cert.ClientCert.Path.Key)
	if err != nil {
		log.Fatal("LoadX509KeyPair err:", err)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost:8080")
	if err != nil {
		log.Fatal("client error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}
