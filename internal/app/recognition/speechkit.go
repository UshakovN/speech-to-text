package recognition

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type YandexResponse struct {
	Result string `json:"result"`
}

type YandexClient struct {
	serviceKey    string
	speechRequest string
	params        string
	response      YandexResponse
}

func NewYandexClient(config *Config) *YandexClient {
	return &YandexClient{
		serviceKey:    config.YandexToken,
		params:        config.RecognitionParams,
		speechRequest: "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize?",
		// strings.Join([]string{"topic=general", "lang=ru-RU"}, "&"),
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
	req, err := http.NewRequest(http.MethodPost, client.speechRequest, speechFile)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("Authorization", "Api-Key "+client.serviceKey)
	respTranscription, err := httpClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer respTranscription.Body.Close()
	bodyTransription, err := ioutil.ReadAll(respTranscription.Body)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(bodyTransription, &client.response)
	if err != nil {
		log.Panic(err)
	}
	return client.response.Result
}
