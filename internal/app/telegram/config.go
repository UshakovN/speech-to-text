package telegram

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessToken string `yaml:"access_token"`
}

func NewConfig() *Config {
	return &Config{
		AccessToken: "",
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
