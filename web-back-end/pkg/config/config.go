package config

import (
	"github.com/BurntSushi/toml"
)

type (
	config struct {
		Port     string
		Postgres struct {
			Dsn string
		}
		Logger struct {
			Level string
		}
		Session struct {
			Name   string
			Path   string
			Secret string
		}
		Applet struct {
			Url                string
			Code2sessionKeyUrl string
			Id                 string
			Secret             string
		}
	}
)

var Conf config

func init() {
	file := "config/config.toml"

	toml.DecodeFile(file, &Conf)
}
