package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type YandexResponse struct {
	Result string `json:"result"`
}

type YandexClient struct {
	ServiceKey    string
	SpeechRequest string
	Params        string
	Response      YandexResponse
}

func NewYandexClient() *YandexClient {
	return &YandexClient{
		ServiceKey:    "AQVNwbhyxMeFNrOcHO6sthjsOgd5gzq3EPhIkCir",
		SpeechRequest: "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize?",
		Params:        strings.Join([]string{"topic=general", "lang=ru-RU"}, "&"),
	}
}

func (client *YandexClient) GetFileChunks(directLink string) (io.Reader, error) {
	file, err := http.Get(directLink)
	if err != nil {
		return nil, err
	}
	return file.Body, nil
}

func (client *YandexClient) GetTranscription(speechFile io.Reader) string {
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(http.MethodPost, client.SpeechRequest, speechFile)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("Authorization", "Api-Key "+client.ServiceKey)
	respTranscription, err := httpClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer respTranscription.Body.Close()
	bodyTransription, err := io.ReadAll(respTranscription.Body)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(bodyTransription, &client.Response)
	if err != nil {
		log.Panic(err)
	}
	return client.Response.Result
}
