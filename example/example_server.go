package main

import (
	"github/j92z/go_kaadda/genx"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type httpsHandler struct {
}

func (*httpsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "golang https server!!!")
}

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

	s := &http.Server{
		Addr:    ":8080",
		Handler: &httpsHandler{},
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	if err = s.ListenAndServeTLS(cert.ServerCert.Path.Ca, cert.ServerCert.Path.Key); err != nil {
		log.Fatal("ListenAndServeTLS err:", err)
	}
}
