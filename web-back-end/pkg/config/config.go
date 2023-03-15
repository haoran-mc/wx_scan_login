package config

import (
	"github.com/BurntSushi/toml"
)

type (
	config struct {
		Address string
		Port    string
		Mysql   struct {
			Dsn    string
			User   string
			Pass   string
			Dbname string
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
			Port               string
			Code2sessionKeyUrl string
			Id                 string
			Secret             string
		}
		Web struct {
			Router struct {
				Scanned string
				Welcome string
			}
		}
	}
)

var Conf config

func init() {
	file := "config/config.toml"
	toml.DecodeFile(file, &Conf)
}
