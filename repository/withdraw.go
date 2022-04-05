package repository

import (
	"dk-project-service/entity"
	"dk-project-service/utils"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type (
	WdRepo interface {
		GetAllWd() ([]entity.WdReqDetail, error)
		GetAllWdInWeek() ([]entity.WdReqDetail, error)
		GetWdById(id string) (entity.WithdrawRequest, error)

		GetAllWdByUserID(userId int) ([]entity.WdReqModel, error)

		GetWdReqInWeekByUserID(userId int) (entity.WdReqModel, error)

		CreateWdReq(data entity.WdReqModel) error
		UpdateWdReqByID(update entity.WdReqModel) error

		ApproveWdReqById(id string, input entity.UpdateWdReqApprove) error
	}

	wdRepo struct {
		db *gorm.DB
	}
)

func NewWdRepo(db *gorm.DB) *wdRepo {
	return &wdRepo{db: db}
}

func (r *wdRepo) GetAllWd() ([]entity.WdReqDetail, error) {
	var query = `SELECT 
		wr.id, 
		wr.money_balance, 
		wr.ro_balance, 
		wr.ro_money_balance,
		wr.created_at,
		wr.updated_at,
		wr.approved,
		u.id as user_id,
		u.fullname,
		u.phone_number,
		ba.id as bank_acc_id,
		ba.bank_name,
		ba.bank_number,
		ba.name_on_bank
	FROM withdraw_requests wr
	JOIN users u ON u.id = wr.user_id
	JOIN bank_accounts ba ON ba.id = wr.bank_acc_id;`

	var wdReqs []entity.WdReqDetail

	if err := r.db.Raw(query).Scan(&wdReqs).Error; err != nil {
		return nil, err
	}

	return wdReqs, nil
}

func (r *wdRepo) GetAllWdInWeek() ([]entity.WdReqDetail, error) {
	var query = `SELECT 
		wr.id, 
		wr.money_balance, 
		wr.ro_balance, 
		wr.ro_money_balance,
		wr.created_at,
		wr.updated_at,
		wr.approved,
		u.id as user_id,
		u.fullname,
		u.phone_number,
		ba.id as bank_acc_id,
		ba.bank_name,
		ba.bank_number,
		ba.name_on_bank
	FROM withdraw_requests wr
	JOIN users u ON u.id = wr.user_id
	JOIN bank_accounts ba ON ba.id = wr.bank_acc_id
	WHERE created_at BETWEEN ? AND ?;`

	var wdReqs []entity.WdReqDetail

	startDateStr, endDateStr, err := utils.CountRangeDate()
	if err != nil {
		return nil, err
	}

	if err := r.db.Raw(query, startDateStr, endDateStr).Scan(&wdReqs).Error; err != nil {
		return nil, err
	}

	return wdReqs, nil
}

func (r *wdRepo) GetWdById(id string) (entity.WithdrawRequest, error) {
	var wdReq entity.WithdrawRequest

	if err := r.db.Where("id = ?", id).Find(&wdReq).Error; err != nil {
		return wdReq, err
	}

	return wdReq, nil
}

func (r *wdRepo) GetWdReqInWeekByUserID(userId int) (entity.WdReqModel, error) {
	var wdGorm entity.WithdrawRequest

	userIdStr := strconv.Itoa(userId)

	startDateStr, endDateStr, err := utils.CountRangeDate()
	if err != nil {
		return entity.WdReqModel{}, err
	}

	if err := r.db.Where("created_at >= ? AND created_at <= ? AND user_id = ?", startDateStr, endDateStr, userIdStr).Find(&wdGorm).Error; err != nil {
		return entity.WdReqModel{}, err
	}

	if wdGorm.Id == 0 && wdGorm.UserId == 0 && wdGorm.BankAccId == 0 {
		return entity.WdReqModel{}, nil
	} else {
		wdReq, err := wdGorm.ToWdReqModel()
		if err != nil {
			return wdReq, err
		}

		return wdReq, nil
	}
}

func (r *wdRepo) GetAllWdByUserID(userId int) ([]entity.WdReqModel, error) {
	var wdGorm []entity.WithdrawRequest
	var wdReqs []entity.WdReqModel

	if err := r.db.Where("user_id = ?", userId).Find(&wdGorm).Error; err != nil {
		return nil, err
	}

	for _, wd := range wdGorm {
		wdReq, err := wd.ToWdReqModel()
		if err != nil {
			return nil, err
		}

		wdReqs = append(wdReqs, wdReq)
	}

	return wdReqs, nil
}

func (r *wdRepo) CreateWdReq(data entity.WdReqModel) error {
	var flag = ""

	var queryField = "user_id, bank_acc_id, created_at, updated_at, approved"
	var queryEmpty = "?, ?, ?, ?, ?"

	if data.MoneyBalance != 0 {
		flag += "m"
		queryField += ", money_balance"
		queryEmpty += ", ?"
	}

	if data.RoBalance != 0 {
		flag += "r"
		queryField += ", ro_balance, ro_money_balance"
		queryEmpty += ", ?, ?"
	}

	var query string = fmt.Sprintf("INSERT INTO withdraw_requests (%s) VALUES (%s)", queryField, queryEmpty)

	switch flag {
	case "m":
		if err := r.db.Exec(query, data.UserId, data.BankAccId, time.Now().Add(7*time.Hour), time.Now().Add(7*time.Hour), false, data.MoneyBalance).Error; err != nil {
			return err
		}
	case "r":
		if err := r.db.Exec(query, data.UserId, data.BankAccId, time.Now().Add(7*time.Hour), time.Now().Add(7*time.Hour), false, data.RoBalance, data.RoMoneyBalance).Error; err != nil {
			return err
		}
	case "mr":
		if err := r.db.Exec(query, data.UserId, data.BankAccId, time.Now().Add(7*time.Hour), time.Now().Add(7*time.Hour), false, data.MoneyBalance, data.RoBalance, data.RoMoneyBalance).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *wdRepo) UpdateWdReqByID(update entity.WdReqModel) error {
	var flag = ""

	var queryField = "updated_at = ?, approved = ?"

	if update.MoneyBalance != 0 {
		flag += "m"
		queryField += ", money_balance = ?"
	}

	if update.RoBalance != 0 {
		flag += "r"
		queryField += ", ro_balance = ?, ro_money_balance = ?"
	}

	var query = fmt.Sprintf("UPDATE withdraw_requests SET %s WHERE id = ?", queryField)

	switch flag {
	case "m":
		if err := r.db.Exec(query, time.Now().Add(7*time.Hour), false, update.MoneyBalance, update.Id).Error; err != nil {
			return err
		}
	case "r":
		if err := r.db.Exec(query, time.Now().Add(7*time.Hour), false, update.RoBalance, update.RoMoneyBalance, update.Id).Error; err != nil {
			return err
		}
	case "mr":
		if err := r.db.Exec(query, time.Now().Add(7*time.Hour), false, update.MoneyBalance, update.RoBalance, update.RoMoneyBalance, update.Id).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *wdRepo) ApproveWdReqById(id string, input entity.UpdateWdReqApprove) error {
	var query = "UPDATE withdraw_requests SET approved = ?, updated_at = ? WHERE id = ?"

	if err := r.db.Exec(query, input.Approved, time.Now(), id).Error; err != nil {
		return err
	}

	return nil
}
