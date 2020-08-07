package httpx

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github/j92z/go_kaadda/genx"
	"github/j92z/go_kaadda/httpx/enginex"
	"github/j92z/go_kaadda/setting"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Run() {
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.EnvSetting.Server.Port),
		Handler:        enginex.InitEngineX(),
		ReadTimeout:    time.Duration(setting.EnvSetting.Server.ReadTimeOut) * time.Second,
		WriteTimeout:   time.Duration(setting.EnvSetting.Server.WriteTimeOut) * time.Second,
		MaxHeaderBytes: setting.EnvSetting.Server.MaxHeaderBytes,
	}
	if setting.EnvSetting.Server.Tls {
		cert := GetCert()
		pool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(cert.RootCert.Path.Ca)
		if err != nil {
			log.Fatal("ReadFile err:", err)
			return
		}
		pool.AppendCertsFromPEM(ca)
		httpServer.TLSConfig = &tls.Config{
			ClientCAs: pool,
		}
		if setting.EnvSetting.Server.ClientAuth {
			httpServer.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
		if err := httpServer.ListenAndServeTLS(cert.ServerCert.Path.Ca, cert.ServerCert.Path.Key); err != nil {
			log.Fatal("ListenAndServeTLS err:", err)
		}
	} else {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("ListenAndServe err:", err)
		}
	}
}

func GetCert() *genx.Cert {
	return genx.New(genx.CertConfig{
		Path: setting.EnvSetting.Certificate.Path,
		Name: setting.EnvSetting.Certificate.Name,
		Host: setting.EnvSetting.Certificate.Hosts,
	})
}
