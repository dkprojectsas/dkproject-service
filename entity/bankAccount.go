package entity

type BankAccount struct {
	Id         int    `json:"id"`
	UserId     int    `json:"user_id"`
	BankName   string `json:"bank_name"`
	BankNumber string `json:"bank_number"`
	NameOnBank string `json:"name_on_bank"`
}

type BankAccountInput struct {
	UserId     int    `json:"user_id" binding:"required"`
	BankName   string `json:"bank_name" binding:"required"`
	BankNumber string `json:"bank_number" binding:"required"`
	NameOnBank string `json:"name_on_bank" binding:"required"`
}
