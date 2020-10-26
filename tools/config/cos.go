package config

import (
	"github.com/spf13/viper"
)

type Cos struct {
	Url   			string
	Secretid   		string
	Secretkey   	string
	Flag   			bool
}

func InitCos(cfg *viper.Viper) *Cos {
	return &Cos{
		Url:			cfg.GetString("url"),
		Secretid:		cfg.GetString("secretid"),
		Secretkey:		cfg.GetString("secretkey"),
		Flag:			cfg.GetBool("flag"),
	}
}

var CosConfig = new(Cos)
