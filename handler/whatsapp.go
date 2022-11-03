package handler

// type (
// 	WhatsappHandler interface {
// 		SendWhatsappMessage(data entity.WhatsappRequestBody)
// 	}

// 	whatsappHandler struct {
// 		client *twilio.RestClient
// 	}
// )

// func NewWhatsappHandler(client *twilio.RestClient) *whatsappHandler {
// 	return &whatsappHandler{client: client}
// }

// func (h *whatsappHandler) SendWhatsappMessage(data entity.WhatsappRequestBody) {
// 	err := godotenv.Load()
// 	config.FailOnError(err, 27, "handler/whatsapp.go")

// 	params := &openapi.CreateMessageParams{}
// 	params.SetTo(data.To)
// 	params.SetFrom(data.From)
// 	params.SetBody(data.Text)

// 	resp, err := h.client.ApiV2010.CreateMessage(params)
// 	if err != nil {
// 		config.FailOnError(err, 40, "handler/whatsapp.go")
// 	} else {
// 		response, _ := json.Marshal(*resp)
// 		fmt.Println("Response: " + string(response))
// 	}
// }
