package recognition

type Config struct {
	TelegramToken     string `yaml:"telegram_token"`
	YandexToken       string `yaml:"yandex_token"`
	RecognitionParams string `yaml:"recognition_params"`
}

func NewConfig() *Config {
	return &Config{
		TelegramToken:     "",
		YandexToken:       "",
		RecognitionParams: "",
	}
}
