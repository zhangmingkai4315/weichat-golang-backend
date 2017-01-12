package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Security struct {
		Token string
	}
	Server struct {
		Host string
		Port int
	}
}

var ConfigObj Config

func NewConfig() (*Config, error) {
	if ConfigObj.Server.Host == "" {
		if _, err := toml.DecodeFile("./config/config.toml", &ConfigObj); err != nil {
			log.Printf("Config Error %s", err)
			return nil, err
		}
		log.Printf("Loading Config File Success")
	}
	return &ConfigObj, nil
}
