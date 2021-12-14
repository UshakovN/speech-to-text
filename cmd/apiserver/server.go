package main

import (
	"flag"
	"log"

	"github.com/UshakovN/speech-to-text/internal/app/apiserver"
)

var (
	pathconfigServer string
)

func init() {
	flag.StringVar(&pathconfigServer, "server-config", "configs/server.yaml", "path to")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	if err := config.UnmarshalYaml(pathconfigServer); err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
