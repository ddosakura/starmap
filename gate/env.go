package main

import (
	"github.com/micro/go-log"
	"github.com/timest/env"
)

type starmap struct {
	Gateway struct {
		Addr string `default:":8080"`
	}
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
