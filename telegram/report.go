package telegram

import (
	"fmt"
	"gitlab.novgate.com/common/common/logger"
	"gitlab.novgate.com/common/common/utils"
	"net/http"
	"net/url"
	"strconv"
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

func PayInError(payChannel string) {

	payChannel = url.QueryEscape(payChannel)
	uri := fmt.Sprintf("/payinError?payChannel=%s", payChannel)
	go do(uri)
}

func Payout(orderNo, payChannel string, money int64) {
	encodedPayChannel := url.QueryEscape(payChannel)
	encodedOrderNo := url.QueryEscape(orderNo)
	encodedMoney := url.QueryEscape(fmt.Sprintf("%0.2f", utils.Cent2Yuan(money)))
	encodedPlatform := url.QueryEscape(telegram.Platform)
	uri := fmt.Sprintf(
		"/payoutWarring?payChannel=%s&orderNo=%s&amount=%s&platform=%s",
		encodedPayChannel,
		encodedOrderNo,
		encodedMoney,
		encodedPlatform,
	)
	go do(uri)
}

func OddsSendGifToChannel(userId int, allWins, jackpotBonus, betAmount int64, gameName string) {
	// 定义参数
	uid := strconv.Itoa(userId)
	winMoney := fmt.Sprintf("%0.2f", float64(allWins+jackpotBonus-betAmount)/100)

	encodedUID := url.QueryEscape(uid)
	encodedGameName := url.QueryEscape(gameName)
	encodedWinMoney := url.QueryEscape(winMoney)
	uri := fmt.Sprintf("/sendgif?uid=%s&gameName=%s&winMoney=%s", encodedUID, encodedGameName, encodedWinMoney)
	go do(uri)
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
