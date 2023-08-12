package main

import (
	"flag"
	"fmt"
	"jwt-auth/config"
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
	fmt.Println(conf)
}
