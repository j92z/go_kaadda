package setting

import (
	"github.com/timest/env"
	"strconv"
)

type Setting struct {
	Server struct {
		Host           string `default:"localhost"`
		Port           int    `default:"8888"`
		Path           string
		RunMode        string `default:"release"`
		Tls            bool   `default:"true"`
		ClientAuth     bool   `default:"false"`
		ReadTimeOut    int    `default:"60"`
		WriteTimeOut   int    `default:"60"`
		MaxHeaderBytes int    `default:"1048576"`
	}
	Certificate struct {
		Path  string   `default:"certificate"`
		Name  string   `default:"kaadda"`
		Hosts []string `default:"localhost,127.0.0.1" slice_sep:","`
	}
	SQLite struct {
		TablePrefix string
		LogMode     bool   `default:"false"`
		Path        string `default:"data.db"`
	}
	File struct {
		Path    string `default:"file/"`
		MaxSize int64  `default:"4194304"`
	}
}

var EnvSetting *Setting

func Setup() {
	EnvSetting = new(Setting)
	if err := env.Fill(EnvSetting); err != nil {
		panic(err)
	}
	schema := "https://"
	if !EnvSetting.Server.Tls {
		schema = "http://"
	}
	EnvSetting.Server.Path = schema + EnvSetting.Server.Host + ":" + strconv.Itoa(EnvSetting.Server.Port)
}
