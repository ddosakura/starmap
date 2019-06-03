package common

import (
	"github.com/micro/go-log"
	"github.com/timest/env"
)

type starmap struct {
	Repo struct {
		Mongo struct {
			URL string `env:"URL" default:"root:starpass@127.0.0.1:27017"`
		} `env:"MONGO"`
		Redis struct {
			URL string `env:"URL" default:"127.0.0.1:6379"`
		} `env:"REDIS"`
		MySQL struct {
			Host string `env:"HOST" default:"127.0.0.1:3306"`
			User string `env:"USER" default:"root"`
			Pass string `env:"PASS" default:"starpass"`
			DB   string `env:"DB" default:"starmap"`
		} `env:"MYSQL"`
	} `env:"REPO"`
	Key struct {
		// JWT string `env:"JWT" require:"true"`
		JWT string `env:"JWT" default:"starmap"`
	} `env:"KEY"`
}

func (c *starmap) KeyJWT() []byte {
	return []byte(c.Key.JWT)
}

// config(s)
var (
	ENV = new(starmap)
)

func init() {
	if err := env.Fill(ENV); err != nil {
		log.Fatal(err)
	}
}
