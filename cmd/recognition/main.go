package main

import (
	"flag"
	"log"

	"github.com/UshakovN/speech-to-text/internal/app/recognition"
	"github.com/UshakovN/speech-to-text/internal/app/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	pathconfigRec string
	pathconfigTg  string
)

func init() {
	flag.StringVar(&pathconfigRec, "yandex-config", "configs/recognition.yaml", "path to")
	flag.StringVar(&pathconfigTg, "telegram-config", "configs/telegram.yaml", "path to")
}

func main() {
	flag.Parse()

	yandexConfig := recognition.NewConfig()
	if err := yandexConfig.UnmarshalYaml(pathconfigRec); err != nil {
		log.Panic(err)
	}

	telegramConfig := telegram.NewConfig()
	if err := telegramConfig.UnmarshalYaml(pathconfigTg); err != nil {
		log.Panic(err)
	}

	telegramBot, err := tgbotapi.NewBotAPI(telegramConfig.AccessToken)
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

		}
		if update.Message.Voice == nil {
			continue
		}
		if update.Message.Voice.Duration != 0 {
			log.Printf("Пользователь [@%s], Длительность [%d]", update.Message.From.UserName, update.Message.Voice.Duration)
			speech, err := telegramBot.GetFileDirectURL(update.Message.Voice.FileID)
			if err != nil {
				log.Panic(err.Error())
			}
			yandexClient := recognition.NewYandexClient(yandexConfig)
			chunks, err := yandexClient.GetFileChunks(speech)
			if err != nil {
				log.Panic(err.Error())
			}
			transcription := yandexClient.GetTranscription(chunks)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, transcription)
			telegramBot.Send(msg)
		}
	}
}
