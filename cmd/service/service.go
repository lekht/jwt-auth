package main

import (
	"flag"
	"jwt-auth/config"
	"jwt-auth/internal/app"
	"log"
)

var conf config.Config

func init() {
	path := flag.String("config", "", "path to config file")
	flag.Parse()

	err := conf.LoadEnv(*path)
	if err != nil {
		log.Fatalf("failed to load environment variables: %s", err)
	}
}

func main() {
	app.Run(&conf)
}
