package app

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string
	Host string
}

func ServerConfigure() ServerConfig {
	const LocalHost = "0.0.0.0"

	viper.SetDefault("Host", LocalHost)
	viper.SetDefault("Port", "5000")

	viper.AutomaticEnv()

	var srvCfg ServerConfig
	if err := viper.Unmarshal(&srvCfg); err != nil {
		log.Panicf("It was impossible configure Server. %+v", err)
	}

	return srvCfg
}
