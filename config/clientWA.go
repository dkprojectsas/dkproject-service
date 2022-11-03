package config

import (
	"github.com/twilio/twilio-go"
)

func ClientWA() *twilio.RestClient {
	accSID, authToken := GetEnvWA()

	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username:   accSID,
		Password:   authToken,
		AccountSid: accSID,
	})

	return client
}
