package apiserver

import (
	"io/ioutil"

	"github.com/UshakovN/speech-to-text/internal/app/store"
	"gopkg.in/yaml.v2"
)

type Config struct {
	BindAddr string `yaml:"bind_addr"`
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		Store:    store.NewConfig(),
	}
}

func (config *Config) UnmarshalYaml(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return err
	}
	return nil
}
