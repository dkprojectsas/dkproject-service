package repository

import (
	"dk-project-service/entity"
	"dk-project-service/utils"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type (
	TransRepo interface {
		InsertTrans(trans entity.TransInput) error
		BulkInsertTrans(trans []entity.TransInput) error
		GetTransactionById(id int) ([]entity.Transaction, error)
		GetByCategory(cat string) ([]entity.Transaction, error)
		GetAllCatForAdmin() ([]entity.Transaction, error)
	}

	transRepo struct {
		db *gorm.DB
	}
)

func NewTransRepo(db *gorm.DB) *transRepo {
	return &transRepo{db: db}
}

var queryBaseTransaction = `
	SELECT 
		t.id, 
		t.from_id, 
		u.fullname as from_fullname,
		u.username as from_username,
		t.to_id,
		u2.fullname as to_fullname,
		u2.username as to_username,
		t.description, 
		t.category, 
		t.sas_balance, 
		t.ro_balance, 
		t.money_balance, 
		t.ro_money_balance, 
		t.created_at 
	FROM transactions t  
	JOIN users u ON u.id = t.from_id
	JOIN users u2 ON u2.id = t.to_id`

func (r *transRepo) InsertTrans(trans entity.TransInput) error {
	var query string

	if trans.SASBalance != 0 {
		query = `INSERT INTO transactions (from_id, to_id, category, description, sas_balance) VALUES (?, ?, ?, ?, ?)`

		if err := r.db.Exec(query, trans.FromId, trans.ToId, trans.Category, trans.Description, trans.SASBalance).Error; err != nil {
			return err
		}
	}

	if trans.ROBalance != 0 {
		query = `INSERT INTO transactions (from_id, to_id, category, description, ro_balance) VALUES (?, ?, ?, ?, ?)`

		if err := r.db.Exec(query, trans.FromId, trans.ToId, trans.Category, trans.Description, trans.ROBalance).Error; err != nil {
			return err
		}
	}

	if trans.MoneyBalance != 0 {
		query = `INSERT INTO transactions (from_id, to_id, category, description, money_balance) VALUES (?, ?, ?, ?, ?)`

		if err := r.db.Exec(query, trans.FromId, trans.ToId, trans.Category, trans.Description, trans.MoneyBalance).Error; err != nil {
			return err
		}
	}

	if trans.ROMoneyBalance != 0 {
		query = `INSERT INTO transactions (from_id, to_id, category, description, ro_money_balance) VALUES (?, ?, ?, ?, ?)`

		if err := r.db.Exec(query, trans.FromId, trans.ToId, trans.Category, trans.Description, trans.ROMoneyBalance).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *transRepo) BulkInsertTrans(trans []entity.TransInput) error {
	var queryInsert = `INSERT INTO transactions (from_id, to_id, category, description, sas_balance, ro_balance, money_balance, ro_money_balance) VALUES `

	for i, t := range trans {
		switch {
		case t.SASBalance != 0:
			queryInsert += fmt.Sprintf(`(%d, %d, '%s', '%s', %d, NULL, NULL, NULL)`, t.FromId, t.ToId, t.Category, t.Description, t.SASBalance)
		case t.ROBalance != 0:
			queryInsert += fmt.Sprintf(`(%d, %d, '%s', '%s', NULL, %d, NULL, NULL)`, t.FromId, t.ToId, t.Category, t.Description, t.ROBalance)
		case t.MoneyBalance != 0:
			queryInsert += fmt.Sprintf(`(%d, %d, '%s', '%s', NULL, NULL, %d, NULL)`, t.FromId, t.ToId, t.Category, t.Description, t.MoneyBalance)
		case t.ROMoneyBalance != 0:
			queryInsert += fmt.Sprintf(`(%d, %d, '%s', '%s', NULL, NULL, NULL, %d)`, t.FromId, t.ToId, t.Category, t.Description, t.ROMoneyBalance)
		}

		if (i + 1) != len(trans) {
			queryInsert += ", "
		}
	}

	if err := r.db.Exec(queryInsert).Error; err != nil {
		return err
	}

	return nil
}

func (r *transRepo) GetTransactionById(id int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction

	idStr := strconv.Itoa(id)

	_, m, y := utils.GetDateNow()

	query := queryBaseTransaction + " WHERE ( from_id = ? OR to_id = ? ) AND MONTH(created_at) = ? AND YEAR(created_at) = ? ORDER BY created_at DESC"

	if err := r.db.Raw(query, idStr, idStr, m, y).Scan(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *transRepo) GetByCategory(cat string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction

	_, m, y := utils.GetDateNow()

	query := queryBaseTransaction + " WHERE category = ? AND MONTH(created_at) = ? AND YEAR(created_at) = ? ORDER BY created_at DESC"

	if err := r.db.Raw(query, cat, m, y).Scan(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *transRepo) GetAllCatForAdmin() ([]entity.Transaction, error) {
	var transactions []entity.Transaction

	query := `
	SELECT 
		t.id, 
		t.from_id, 
		u.fullname as from_fullname,
		u.username as from_username,
		t.to_id,
		u2.fullname as to_fullname,
		u2.username as to_username,
		t.description, 
		t.category, 
		t.sas_balance, 
		t.ro_balance, 
		t.money_balance, 
		t.ro_money_balance, 
		t.created_at 
	FROM transactions t  
	JOIN users u ON u.id = t.from_id
	JOIN users u2 ON u2.id = t.to_id 
	WHERE (t.from_id = 1 OR t.to_id = 1)
	 	AND ( t.category IN ('admin_fee', 'sas_balance', 'ro_balance') OR (t.category = 'umum' AND t.description NOT LIKE 'bonus %'))
		AND MONTH(created_at) = ? 
		AND YEAR(created_at) = ?
	ORDER BY t.created_at DESC
	`

	_, m, y := utils.GetDateNow()
	if err := r.db.Raw(query, m, y).Scan(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}
