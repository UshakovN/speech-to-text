package store

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DatabaseUrl string `yaml:"database_url"`
}

func NewConfig() *Config {
	return &Config{
		DatabaseUrl: "",
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
