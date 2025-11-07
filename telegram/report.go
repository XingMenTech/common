package telegram

import (
	"fmt"
	"github.com/XingMenTech/common/logger"
	"github.com/google/go-querystring/query"
	"net/http"
)

var telegram *Config

type Config struct {
	Api      string `yaml:"api"`
	Key      string `yaml:"key"`
	Platform string `yaml:"platform"`
}

func InitTelegram(config *Config) {
	if config == nil {
		return
	}

	telegram = config
}

func SendMessage(uri string, param map[string]interface{}) {
	values, err := query.Values(param)
	if err != nil {
		return
	}
	go do(uri + "?" + values.Encode())
}

func do(uri string) {
	reqUrl := fmt.Sprintf("%s%s", telegram.Api, uri)
	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		logger.LOG.Warn("sent to telegram bot service errors: ", err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KEY", telegram.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.LOG.Warn("Failed to send request: ", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.LOG.Warn("Request failed with status: ", resp.Status)
	}
}
