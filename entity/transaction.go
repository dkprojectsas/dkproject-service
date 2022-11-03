package entity

type Transaction struct {
	Id             int    `json:"id"`
	FromId         int    `json:"from_id"`
	FromFullname   string `json:"from_fullname"`
	FromUsername   string `json:"from_username"`
	ToId           int    `json:"to_id"`
	ToFullname     string `json:"to_fullname"`
	ToUsername     string `json:"to_username"`
	Description    string `json:"description"`
	Category       string `json:"category"`
	SASBalance     int    `json:"sas_balance"`
	ROBalance      int    `json:"ro_balance"`
	MoneyBalance   int    `json:"money_balance"`
	ROMoneyBalance int    `json:"ro_money_balance"`
	CreatedAt      string `json:"created_at"`
}

type TransInput struct {
	FromId         int    `json:"from_id" binding:"required"`
	ToId           int    `json:"to_id" binding:"required"`
	Description    string `json:"description" binding:"required"`
	Category       string `json:"category" binding:"required"`
	SASBalance     int    `json:"sas_balance"`
	ROBalance      int    `json:"ro_balance"`
	MoneyBalance   int    `json:"money_balance"`
	ROMoneyBalance int    `json:"ro_money_balance"`
}

type NewDownlineInput struct {
	UplineId     int `json:"upline_id"`
	MoneyBalance int `json:"money_balance"`
}

type BuySASAdminInput struct {
	UserId       int `json:"user_id"`
	SASBalance   int `json:"sas_balance"`
	MoneyBalance int `json:"money_balance"`
}
type BuyROAdminInput struct {
	UserId       int `json:"user_id"`
	ROBalance    int `json:"ro_balance"`
	MoneyBalance int `json:"money_balance"`
}

type AddBalanceInput struct {
	ROBalance  int `json:"ro_balance"`
	SASBalance int `json:"sas_balance"`
}

const (
	TransCategoryGeneral  = "umum"
	TransCategorySAS      = "sas_balance"
	TransCategoryRO       = "ro_balance"
	TransCategoryAdminFee = "admin_fee"
)
