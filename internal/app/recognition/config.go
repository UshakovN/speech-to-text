package recognition

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	YandexToken       string `yaml:"yandex_token"`
	RecognitionParams string `yaml:"recognition_params"`
}

func NewConfig() *Config {
	return &Config{
		YandexToken:       "",
		RecognitionParams: "",
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
