package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type config struct {
	Security struct {
		Token     string
		AppId     string
		AppSecret string
	}
	Server struct {
		Host string
		Port int
	}
	Database struct {
		Host     string
		User     string
		Password string
		Port     int
		DBName   string
	}
	Url struct {
		TokenRefreshUrl string
	}
}

var ConfigObj config

func init() {
	NewConfig(&ConfigObj)
}
func NewConfig(c *config) (*config, error) {
	if c.Server.Host == "" {
		if _, err := toml.DecodeFile("./config/config.toml", &c); err != nil {
			log.Printf("Config Error %s", err)
			return nil, err
		}
		log.Printf("Loading Config File Success")
	}
	return c, nil
}
