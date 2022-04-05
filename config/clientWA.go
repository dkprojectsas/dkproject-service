package config

import (
	"os"

	"github.com/twilio/twilio-go"
)

func ClientWA() *twilio.RestClient {
	accSID := os.Getenv("T_ACC_SID")
	authToken := os.Getenv("T_AUTH_TOKEN")

	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username:   accSID,
		Password:   authToken,
		AccountSid: accSID,
	})

	return client
}
