package handler

import (
	"bytes"
	"dk-project-service/config"
	"dk-project-service/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type (
	SMSHandler interface {
		SendSMS(message string)
	}

	smsHandler struct {
	}
)

func NewSMSHandler() *smsHandler {
	return &smsHandler{}
}

func (h *smsHandler) SendSMS(sendTo string, message string) {
	err := godotenv.Load()
	config.FailOnError(err, 32, "handler/sms.go")

	body := entity.SMSRequestBody{
		APIKey:    os.Getenv("V_API_KEY"),
		APISecret: os.Getenv("V_API_SECRET"),
		To:        "+" + sendTo,
		From:      "DK Admin",
		Text:      message,
	}

	smsBody, err := json.Marshal(body)
	config.FailOnError(err, 42, "handler/sms.go")

	resp, err := http.Post("https://rest.nexmo.com/sms/json", "application/json", bytes.NewBuffer(smsBody))
	config.FailOnError(err, 45, "handler/sms.go")

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	config.FailOnError(err, 50, "handler/sms.go")

	fmt.Println(string(respBody))
}
