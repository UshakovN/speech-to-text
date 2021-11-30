package main

import (
	"io/ioutil"
	"log"

	"github.com/UshakovN/speech-to-text/internal/app/recognition"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v2"
)

func loadConfig() *recognition.Config {
	config := recognition.NewConfig()

	file, err := ioutil.ReadFile("./configs/main.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func main() {
	config := loadConfig()

	telegramBot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Авторизация [@%s]", telegramBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := telegramBot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Text != "" {
			log.Printf("Пользователь [@%s],  Сообщение [%s]", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			telegramBot.Send(msg)

		} else if update.Message.Voice.Duration != 0 {
			log.Printf("Пользователь [@%s], Длительность [%d]", update.Message.From.UserName, update.Message.Voice.Duration)
			speech, err := telegramBot.GetFileDirectURL(update.Message.Voice.FileID)
			if err != nil {
				log.Panic(err)
			}
			yandexClient := recognition.NewYandexClient(config.YandexToken, config.RecognitionParams)
			chunks, err := yandexClient.GetFileChunks(speech)
			if err != nil {
				log.Panic(err)
			}
			transcription := yandexClient.GetTranscription(chunks)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, transcription)
			telegramBot.Send(msg)
		}
	}
}
