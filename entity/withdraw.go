package entity

import "time"

type WithdrawRequest struct {
	Id             int
	UserId         int
	BankAccId      int
	MoneyBalance   int
	RoBalance      int
	RoMoneyBalance int
	CreatedAt      string
	UpdatedAt      string
	Approved       bool
}

type WdReqDetail struct {
	Id             int    `json:"id"`
	MoneyBalance   int    `json:"money_balance"`
	RoBalance      int    `json:"ro_balance"`
	RoMoneyBalance int    `json:"ro_money_balance"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"update_at"`
	Approved       bool   `json:"approved"`
	// user
	UserId      int    `json:"user_id"`
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phone_number"`
	// bank account
	BankAccId  int    `json:"bank_acc_id"`
	BankName   string `json:"bank_name"`
	BankNumber string `json:"bank_number"`
	NameOnBank string `json:"name_on_bank"`
}

type WdReqModel struct {
	Id             int       `json:"id"`
	UserId         int       `json:"user_id"`
	BankAccId      int       `json:"bank_acc_id"`
	MoneyBalance   int       `json:"money_balance"`
	RoBalance      int       `json:"ro_balance"`
	RoMoneyBalance int       `json:"ro_money_balance"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"update_at"`
	Approved       bool      `json:"approved"`
}

type WdReqInput struct {
	UserId       int `json:"user_id" binding:"required"`
	BankAccId    int `json:"bank_acc_id" binding:"required"`
	Moneybalance int `json:"money_balance"`
	RoBalance    int `json:"ro_balance"`
}

type UpdateWdReqApprove struct {
	Approved bool `json:"approved" binding:"required"`
}

const (
	BonusUser     = 2000
	BonusJaringan = 3000
	BiayaAdmin    = 300
)

const ParseFormat = "2006-01-02 15:04:05"

func (wr *WithdrawRequest) ToWdReqModel() (WdReqModel, error) {
	var wdReq = WdReqModel{
		Id:             wr.Id,
		UserId:         wr.UserId,
		BankAccId:      wr.BankAccId,
		MoneyBalance:   wr.MoneyBalance,
		RoBalance:      wr.RoBalance,
		RoMoneyBalance: wr.RoMoneyBalance,
		Approved:       wr.Approved,
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return wdReq, err
	}

	createdAt, err := time.Parse(ParseFormat, wr.CreatedAt)
	if err != nil {
		return wdReq, err
	}

	updatedAt, err := time.Parse(ParseFormat, wr.UpdatedAt)
	if err != nil {
		return wdReq, err
	}

	localCAt := createdAt.In(loc)
	localUAt := updatedAt.In(loc)

	wdReq.CreatedAt = localCAt
	wdReq.UpdatedAt = localUAt

	return wdReq, nil
}
