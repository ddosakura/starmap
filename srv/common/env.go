package common

import (
	"github.com/micro/go-log"
	"github.com/timest/env"
)

type config struct {
	Repo struct {
		Mongo struct {
			URL string `env:"URL" default:"root:starpass@127.0.0.1:27017"`
		} `env:"MONGO"`
		Redis struct {
			URL string `env:"URL" default:"127.0.0.1:6379"`
		} `env:"REDIS"`
		MySQL struct {
			URL  string `env:"URL" default:"127.0.0.1:3306"`
			User string `env:"USER" default:"root"`
			Pass string `env:"PASS" default:"starpass"`
		} `env:"MYSQL"`
	} `env:"REPO"`
}

// config(s)
var (
	ENV = new(config)
)

func init() {
	if err := env.Fill(ENV); err != nil {
		log.Fatal(err)
	}
}
