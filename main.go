package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	telegramBot, err := tgbotapi.NewBotAPI("2056045746:AAEHVepiBuHuBHTSmN-kBlGDaSDCBbEMWmk")
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
			yandexClient := NewYandexClient()
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
