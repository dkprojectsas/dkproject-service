package service

import (
	"dk-project-service/entity"
	"dk-project-service/repository"
)

type (
	BankAccountService interface {
		Insert(input entity.BankAccountInput) error
		GetByUser(input int) (entity.BankAccount, error)

		UpdateByID(id int, input entity.BankAccountInput) error
	}

	bankAccountService struct {
		baRepo repository.BankAccountRepo
	}
)

func NewBankAccountService(baRepo repository.BankAccountRepo) *bankAccountService {
	return &bankAccountService{baRepo: baRepo}
}

func (s *bankAccountService) Insert(input entity.BankAccountInput) error {
	return s.baRepo.Insert(input)
}

func (s *bankAccountService) GetByUser(input int) (entity.BankAccount, error) {
	return s.baRepo.GetByUserId(input)
}

func (s *bankAccountService) UpdateByID(id int, input entity.BankAccountInput) error {
	bankAccUpdate := entity.BankAccount{
		Id:         id,
		UserId:     input.UserId,
		BankName:   input.BankName,
		BankNumber: input.BankNumber,
		NameOnBank: input.NameOnBank,
	}

	return s.baRepo.UpdateById(bankAccUpdate)
}
