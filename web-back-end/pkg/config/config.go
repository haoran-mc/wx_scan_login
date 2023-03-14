package config

import (
	"github.com/BurntSushi/toml"
)

type (
	config struct {
		Address string `toml:"address"`
		Port    string `toml:"port"`
		Mysql   struct {
			Dsn    string `toml:"dsn"`
			User   string `toml:"user"`
			Pass   string `toml:"pass"`
			Dbname string `toml:"dbname"`
		}
		Logger struct {
			Level string `toml:"level"`
		}
		Session struct {
			Name   string `toml:"name"`
			Path   string `toml:"path"`
			Secret string `toml:"secret"`
		}
		Applet struct {
			Url                string `toml:"url"`
			Code2sessionKeyUrl string `toml:"code2sessionKeyUrl"`
			Id                 string `toml:"id"`
			Secret             string `toml:"secret"`
		}
	}
)

var Conf config

func init() {
	file := "config/config.toml"
	toml.DecodeFile(file, &Conf)
}
