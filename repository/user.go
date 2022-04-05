package repository

import (
	"bytes"
	"dk-project-service/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetuserId(id int) (entity.User, error)

	GetAllUsers() ([]entity.User, error)
	GetUserViews(id string) ([]entity.UserView, error)

	CheckUserLogin(username string, pass string) (entity.User, error)

	GetUsersByParentId(parentId string) ([]entity.User, error)

	// for register repo
	CheckUserId(id int) ([]entity.User, error)
	CreateUser(user entity.User) (entity.User, error)
	UpdateUsername(user entity.User) error
	UpdateUserById(user entity.User) error

	// for transaction
	UpdateBalance(user entity.User) error

	// send WA message credential
	SendWANotification(user entity.User) (entity.WASendResponse, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetuserId(id int) (entity.User, error) {
	var user entity.User

	idStr := strconv.Itoa(id)

	if err := r.db.Where("id = ?", idStr).Find(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) GetAllUsers() ([]entity.User, error) {
	var users []entity.User

	if err := r.db.Where("role = ?", "user").Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (r *userRepository) GetUserViews(id string) ([]entity.UserView, error) {
	var userViews []entity.UserView
	var users []entity.User

	if err := r.db.Not("id = ?", id).Find(&users).Error; err != nil {
		return userViews, err
	}

	for _, u := range users {
		if u.Role == "user" {
			userViews = append(userViews, u.ToUserView())
		}
	}

	return userViews, nil
}

func (r *userRepository) CheckUserId(id int) ([]entity.User, error) {
	var usersDownline []entity.User

	idStr := strconv.Itoa(id)

	if err := r.db.Raw("SELECT * FROM users WHERE parent_id = ? AND position IN ('left', 'right', 'center')", idStr).Scan(&usersDownline).Error; err != nil {
		return usersDownline, err
	}

	return usersDownline, nil
}

func (r *userRepository) CheckUserLogin(username string, pass string) (entity.User, error) {
	var user entity.User

	if err := r.db.Where("username = ? AND password = ?", username, pass).Find(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) CreateUser(user entity.User) (entity.User, error) {
	var query = `INSERT INTO users (id_generate, role, fullname, phone_number, username, password, parent_id, position) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	if err := r.db.Exec(query, user.IdGenerate, user.Role, user.Fullname, user.PhoneNumber, user.Username, user.Password, user.ParentId, user.Position).Error; err != nil {
		return user, err
	}

	if err := r.db.Where("id_generate = ?", user.IdGenerate).Find(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) UpdateUsername(user entity.User) error {
	if err := r.db.Exec("UPDATE users SET username = ? WHERE id_generate = ?", user.Username, user.IdGenerate).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateBalance(user entity.User) error {
	if err := r.db.Exec("UPDATE users SET sas_balance = ?, ro_balance = ?, money_balance = ? WHERE id = ?", user.SASBalance, user.ROBalance, user.MoneyBalance, user.Id).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserById(user entity.User) error {
	log.Println("lakukan edit")

	if err := r.db.Exec("UPDATE users SET username = ?, password = ?, fullname = ?, phone_number = ? WHERE id = ?", user.Username, user.Password, user.Fullname, user.PhoneNumber, user.Id).Error; err != nil {
		return err
	}

	log.Println("success edit")

	return nil
}

func (r *userRepository) GetUsersByParentId(parentId string) ([]entity.User, error) {
	var users []entity.User

	if err := r.db.Where("parent_id = ? && role = ?", parentId, "user").Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (r *userRepository) SendWANotification(user entity.User) (entity.WASendResponse, error) {
	var cbResp entity.WASendResponse

	// err := godotenv.Load()
	// if err != nil {
	// 	return cbResp, err
	// }

	msgReq := fmt.Sprintf("Selamat bergabung di DK, berikut adalah username dan pin anda. Username : %s, PIN / Password : %s \n\nCatatan: ini adalah data rahasia, mohon dijaga baik baik", user.Username, user.Password)

	reqBody := entity.SendWABody{
		UserKey: os.Getenv("ZENZIVA_USER_KEY"),
		PassKey: os.Getenv("ZENZIVA_PASS_KEY"),
		To:      user.PhoneNumber,
		Message: msgReq,
	}

	jsonReq, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://console.zenziva.net/wareguler/api/sendWA/", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		return cbResp, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return cbResp, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cbResp, err
	}

	json.Unmarshal(body, &cbResp)

	log.Println(string(body))

	return cbResp, nil
}
