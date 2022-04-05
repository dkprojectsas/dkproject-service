package repository

import (
	"dk-project-service/entity"
	"strconv"

	"gorm.io/gorm"
)

type (
	BankAccountRepo interface {
		Insert(ba entity.BankAccountInput) error
		GetAll() ([]entity.BankAccount, error)
		GetByUserId(id int) (entity.BankAccount, error)

		UpdateById(update entity.BankAccount) error
	}

	bankAccountRepo struct {
		db *gorm.DB
	}
)

func NewBankAccountRepo(db *gorm.DB) *bankAccountRepo {
	return &bankAccountRepo{db: db}
}

func (r *bankAccountRepo) Insert(ba entity.BankAccountInput) error {

	var query = `INSERT INTO bank_accounts (bank_name, bank_number, name_on_bank, user_id) VALUES (?, ?, ?, ?)`

	if err := r.db.Exec(query, ba.BankName, ba.BankNumber, ba.NameOnBank, ba.UserId).Error; err != nil {
		return err
	}

	return nil
}

func (r *bankAccountRepo) GetByUserId(id int) (entity.BankAccount, error) {
	var bankAccount entity.BankAccount

	idStr := strconv.Itoa(id)

	if err := r.db.Where("user_id = ?", idStr).Find(&bankAccount).Error; err != nil {
		return bankAccount, err
	}

	return bankAccount, nil
}

func (r *bankAccountRepo) GetAll() ([]entity.BankAccount, error) {
	var bankAccounts []entity.BankAccount

	if err := r.db.Find(&bankAccounts).Error; err != nil {
		return bankAccounts, err
	}

	return bankAccounts, nil
}

func (r *bankAccountRepo) UpdateById(update entity.BankAccount) error {
	var query = `UPDATE bank_accounts SET bank_name = ? , bank_number = ? , name_on_bank = ? WHERE id = ?`

	if err := r.db.Exec(query, update.BankName, update.BankNumber, update.NameOnBank, update.Id).Error; err != nil {
		return err
	}

	return nil
}
