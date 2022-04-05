package entity

type SMSRequestBody struct {
	From      string `json:"from"`
	Text      string `json:"text"`
	To        string `json:"to"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

type WhatsappRequestBody struct {
	From string `json:"from"`
	Text string `json:"text"`
	To   string `json:"to"`
}

type WASendResponse struct {
	MessageID int    `json:"messageId"`
	To        string `json:"to"`
	Status    string `json:"status"`
	Text      string `json:"text"`
	Cost      int    `json:"cost"`
}
