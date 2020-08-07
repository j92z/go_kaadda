package genx

import (
	"github/j92z/go_kaadda/pkg/util/file_util"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	rand2 "math/rand"
	"net"
	"strings"
	"time"
)

type Cert struct {
	Config     CertConfig
	RootCert   CertInfo //根文件信息
	ServerCert CertInfo //server文件信息
	ClientCert CertInfo //client文件信息
}

type CertConfig struct {
	Path       string        //文件存储路径
	Name       string        //服务名称
	Host       []string      //host列表
	CertSuffix string        //cert文件后缀
	KeySuffix  string        //key文件后缀
	RsaKeyBits int           //rsa 计算位宽
	ExpireTime time.Duration //证书过期时间
}

type CertInfo struct {
	Certificate *x509.Certificate //根文件
	PrivateKey  *rsa.PrivateKey   //server文件
	Path        CaPath            //client文件
}

type CaPath struct {
	Ca  string //ca 文件路径
	Key string //key 文件路径
}

func New(config ...CertConfig) *Cert {
	if len(config) > 1 {
		panic(errors.New("config length illegal."))
	}
	cert := &Cert{}
	setting := CertConfig{}
	if len(config) == 1 {
		setting = config[0]
	}
	cert.InitConfig(setting)
	cert.InitCert()
	return cert
}

//初始化cert配置
func (c *Cert) InitConfig(setting CertConfig) {
	defaultSetting := CertConfig{
		Path:       defaultPath,
		Name:       defaultName,
		CertSuffix: defaultCertSuffix,
		KeySuffix:  defaultKeySuffix,
		RsaKeyBits: defaultRsaKeyBits,
		ExpireTime: defaultExpireTime,
	}
	if setting.Path != "" {
		if strings.HasSuffix(setting.Path, "/") {
			defaultSetting.Path = setting.Path
		} else {
			defaultSetting.Path = setting.Path + "/"
		}
	}
	if setting.KeySuffix != "" {
		if strings.HasPrefix(setting.KeySuffix, ".") {
			defaultSetting.KeySuffix = setting.KeySuffix
		} else {
			defaultSetting.KeySuffix = "." + setting.KeySuffix
		}
	}
	if setting.CertSuffix != "" {
		if strings.HasPrefix(setting.CertSuffix, ".") {
			defaultSetting.CertSuffix = setting.CertSuffix
		} else {
			defaultSetting.CertSuffix = "." + setting.CertSuffix
		}
	}
	if setting.Name != "" {
		defaultSetting.Name = setting.Name
	}
	if len(setting.Host) > 0 {
		defaultSetting.Host = setting.Host
	}
	if setting.RsaKeyBits > 0 && setting.RsaKeyBits%1024 == 0 {
		defaultSetting.RsaKeyBits = setting.RsaKeyBits
	}
	if setting.ExpireTime > 0 {
		defaultSetting.ExpireTime = setting.ExpireTime
	}
	c.Config = defaultSetting
}

func (c *Cert) InitCert() {
	c.InitContainerFolder()
	c.InitRootCa()
	c.InitServerCa()
	c.InitClientCa()
}

const defaultPath string = "cert/"
const defaultName string = "ca"
const defaultCertSuffix string = ".pem"
const defaultKeySuffix string = ".key"
const defaultRsaKeyBits int = 2048
const defaultCaName string = "ca"
const defaultServerName string = "server"
const defaultClientName string = "client"
const defaultExpireTime time.Duration = 365 * 24 * time.Hour

func (c *Cert) InitRootCa() {

	c.RootCert.Path.Ca = file_util.PathJoin(c.Config.Path, c.Config.Name, defaultCaName+c.Config.CertSuffix)
	c.RootCert.Path.Key = file_util.PathJoin(c.Config.Path, c.Config.Name, defaultCaName+c.Config.KeySuffix)

	if !file_util.CheckFile(c.RootCert.Path.Ca) || !file_util.CheckFile(c.RootCert.Path.Key) {
		c.CreateRootCert()
	}
	certificate, privateKey, err := c.GetCAInfo(c.RootCert.Path.Ca, c.RootCert.Path.Key)
	if err != nil {
		panic(err)
	}

	c.RootCert.Certificate = certificate
	c.RootCert.PrivateKey = privateKey
}

func (c *Cert) InitServerCa() {

	c.ServerCert.Path.Ca = file_util.PathJoin(c.Config.Path, c.Config.Name, defaultServerName+c.Config.CertSuffix)
	c.ServerCert.Path.Key = file_util.PathJoin(c.Config.Path, c.Config.Name, defaultServerName+c.Config.KeySuffix)
	if !file_util.CheckFile(c.ServerCert.Path.Ca) || !file_util.CheckFile(c.ServerCert.Path.Key) {
		c.CreateServerCert()
	}
	certificate, privateKey, err := c.GetCAInfo(c.ServerCert.Path.Ca, c.ServerCert.Path.Key)
	if err != nil {
		panic(err)
	}
	c.ServerCert.Certificate = certificate
	c.ServerCert.PrivateKey = privateKey
}

func (c *Cert) InitClientCa() {

	c.ClientCert.Path.Ca = file_util.PathJoin(c.Config.Path, c.Config.Name, defaultClientName+c.Config.CertSuffix)
	c.ClientCert.Path.Key = file_util.PathJoin(c.Config.Path, c.Config.Name, defaultClientName+c.Config.KeySuffix)
	if !file_util.CheckFile(c.ClientCert.Path.Ca) || !file_util.CheckFile(c.ClientCert.Path.Key) {
		c.CreateClientCert()
	}
	certificate, privateKey, err := c.GetCAInfo(c.ClientCert.Path.Ca, c.ClientCert.Path.Key)
	if err != nil {
		panic(err)
	}
	c.ClientCert.Certificate = certificate
	c.ClientCert.PrivateKey = privateKey
}

func (c *Cert) InitContainerFolder() {
	file_util.CheckDir(file_util.PathJoin(c.Config.Path, c.Config.Name))
}

func (c *Cert) CreateRootCert() {
	ca := c.GenCaStruct(c.Config.Name, true)
	privateKey, _ := rsa.GenerateKey(rand.Reader, c.Config.RsaKeyBits)

	if err := c.CreateCertificate(c.RootCert.Path.Ca, ca, privateKey, ca, nil); err != nil {
		panic(err)
	}
	if err := c.CreatePrivateKey(c.RootCert.Path.Key, privateKey); err != nil {
		panic(err)
	}

}

func (c *Cert) CreateServerCert() {
	name := defaultServerName
	rootCa, rootKey, err := c.GetCAInfo(c.RootCert.Path.Ca, c.RootCert.Path.Key)
	if err != nil {
		panic(err)
	}
	serverCa := c.GenCaStruct(name, false)
	for _, h := range c.Config.Host {
		if ip := net.ParseIP(h); ip != nil {
			serverCa.IPAddresses = append(serverCa.IPAddresses, ip)
		} else {
			serverCa.DNSNames = append(serverCa.DNSNames, h)
		}
	}
	privateKey, _ := rsa.GenerateKey(rand.Reader, c.Config.RsaKeyBits)
	if err := c.CreateCertificate(c.ServerCert.Path.Ca, serverCa, privateKey, rootCa, rootKey); err != nil {
		panic(err)
	}
	if err := c.CreatePrivateKey(c.ServerCert.Path.Key, privateKey); err != nil {
		panic(err)
	}
}

func (c *Cert) CreateClientCert() {
	name := defaultClientName
	rootCa, rootKey, err := c.GetCAInfo(c.RootCert.Path.Ca, c.RootCert.Path.Key)
	if err != nil {
		panic(err)
	}
	clientCa := c.GenCaStruct(name, false)
	privateKey, _ := rsa.GenerateKey(rand.Reader, c.Config.RsaKeyBits)
	if err := c.CreateCertificate(c.ClientCert.Path.Ca, clientCa, privateKey, rootCa, rootKey); err != nil {
		panic(err)
	}
	if err := c.CreatePrivateKey(c.ClientCert.Path.Key, privateKey); err != nil {
		panic(err)
	}
}

//获取root CA证书
func (c *Cert) GetCAInfo(caPath, keyPath string) (*x509.Certificate, *rsa.PrivateKey, error) {
	//解析根证书
	caFile, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, nil, err
	}
	caBlock, _ := pem.Decode(caFile)

	cert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	//解析私钥
	keyFile, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, nil, err
	}
	keyBlock, _ := pem.Decode(keyFile)
	praKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return cert, praKey, nil
}

func (c *Cert) GenCaStruct(orgName string, isCa bool) *x509.Certificate {
	var nowTime = time.Now()
	var expireTime = time.Now().Add(c.Config.ExpireTime)
	return &x509.Certificate{
		SerialNumber: big.NewInt(rand2.Int63()), //证书序列号
		Subject: pkix.Name{
			Country:      []string{"CN"},
			Organization: []string{orgName},
			Province:     []string{"ShangHai"},
			Locality:     []string{"ShangHai"},
		},
		NotBefore:             nowTime,    //证书有效期开始时间
		NotAfter:              expireTime, //证书有效期结束时间
		BasicConstraintsValid: true,       //基本的有效性约束
		IsCA:                  isCa,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
}

func (c *Cert) CreateCertificate(caPath string, cert *x509.Certificate, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey) error {
	publicKey := &key.PublicKey
	privateKey := key
	if caKey != nil {
		privateKey = caKey
	}
	certificateByte, err := x509.CreateCertificate(rand.Reader, cert, caCert, publicKey, privateKey)
	if err != nil {
		return err
	}
	if err := c.WriteCertFile(caPath, certificateByte, Certificate); err != nil {
		return err
	}
	return nil
}

func (c *Cert) CreatePrivateKey(keyPath string, key *rsa.PrivateKey) error {
	privateKeyByte := x509.MarshalPKCS1PrivateKey(key)
	if err := c.WriteCertFile(keyPath, privateKeyByte, RsaPrivateKey); err != nil {
		return err
	}
	return nil
}

type CertFileType int

const (
	_ CertFileType = iota
	Certificate
	PrivateKey
	RsaPrivateKey
)

func (c CertFileType) Name() string {
	switch c {
	case Certificate:
		return "CERTIFICATE"
	case PrivateKey:
		return "PRIVATE KEY"
	case RsaPrivateKey:
		return "RSA PRIVATE KEY"
	}
	return "CERTIFICATE"
}

func (c *Cert) WriteCertFile(filePath string, byteData []byte, fileType CertFileType) error {
	var pemBlock = &pem.Block{
		Type:  fileType.Name(),
		Bytes: byteData,
	}
	byteData64 := pem.EncodeToMemory(pemBlock)
	if err := ioutil.WriteFile(filePath, byteData64, 0777); err != nil {
		return err
	}
	return nil
}
