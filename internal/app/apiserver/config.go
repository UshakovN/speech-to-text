package apiserver

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BindAddr string `yaml:"bind_addr"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
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
