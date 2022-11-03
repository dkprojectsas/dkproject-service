package entity

type User struct {
	Id           int    `json:"id"`
	IdGenerate   string `json:"id_generate"`
	Role         string `json:"role"`
	Fullname     string `json:"fullname"`
	PhoneNumber  string `json:"phone_number"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ParentId     int    `json:"parent_id"`
	Position     string `json:"position"`
	SASBalance   int    `json:"sas_balance"`
	ROBalance    int    `json:"ro_balance"`
	MoneyBalance int    `json:"money_balance"`
}

type UserView struct {
	Id       int    `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	ParentId int    `json:"parent_id"`
	Position string `json:"position"`
}

func (u *User) ToUserView() UserView {
	return UserView{
		Id:       u.Id,
		Fullname: u.Fullname,
		Username: u.Username,
		ParentId: u.ParentId,
		Position: u.Position,
	}
}

type UserId struct {
	Id int `json:"id"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type InputForgotPass struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type InputChangePass struct {
	OldPassword string `json:"old_password" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type ForgotPassResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegister struct {
	Fullname    string `json:"fullname" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	ParentId    int    `json:"parent_id" binding:"required"`
	Position    string `json:"position" binding:"required"`
}

type UserDetail struct {
	Id           int    `json:"id"`
	IdGenerate   string `json:"id_generate"`
	Role         string `json:"role"`
	Fullname     string `json:"fullname"`
	PhoneNumber  string `json:"phone_number"`
	Username     string `json:"username"`
	ParentId     int    `json:"parent_id"`
	Position     string `json:"position"`
	SASBalance   int    `json:"sas_balance"`
	ROBalance    int    `json:"ro_balance"`
	MoneyBalance int    `json:"money_balance"`
}

type UserLoginResponse struct {
	Id          int    `json:"id"`
	Role        string `json:"role"`
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phone_number"`
	ParentId    int    `json:"parent_id"`
	Token       string `json:"auth_token"`
}

type UserUpdateInput struct {
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phone_number"`
}

type SendWABody struct {
	UserKey string `json:"userkey"`
	PassKey string `json:"passkey"`
	To      string `json:"to"`
	Message string `json:"message"`
}
