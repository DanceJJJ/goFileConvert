package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

var cfgDatabase *viper.Viper
var cfgRedis *viper.Viper
var cfgCos *viper.Viper
var cfgCorpwx *viper.Viper
var cfgApplication *viper.Viper
var cfgJwt *viper.Viper
var cfgLog *viper.Viper
var cfgSsl *viper.Viper


//载入配置文件
func ConfigSetup(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	cfgApplication = viper.Sub("settings.application")
	if cfgApplication == nil {
		panic("config not found settings.application")
	}
	ApplicationConfig = InitApplication(cfgApplication)

	cfgCos = viper.Sub("settings.cos")
	if cfgCos == nil {
		panic("config not found settings.cos")
	}
	CosConfig = InitCos(cfgCos)
}

func SetConfig(configPath string, key string, value interface{}) {
	viper.AddConfigPath(configPath)
	viper.Set(key, value)
	viper.WriteConfig()
}
