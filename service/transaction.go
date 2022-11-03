package service

import (
	"dk-project-service/entity"
	"dk-project-service/repository"
	"dk-project-service/utils"
	"errors"
	"fmt"
)

type (
	TransService interface {
		NewRecord(input entity.TransInput) error
		TransactionByUser(id int) ([]entity.Transaction, error)

		InsertNewTrans(input entity.TransInput) error
		InsertBulkTrans(input []entity.TransInput) error

		NewDownline(inputUplineId int) error

		BuySASAdmin(input entity.BuySASAdminInput) error
		BuyROAdmin(input entity.BuyROAdminInput) error

		AddBalanceAdmin(input entity.AddBalanceInput) error

		GetByCategory(cat string) ([]entity.Transaction, error)
		GetAllCategoryForAdmin() ([]entity.Transaction, error)
	}

	transService struct {
		transRepo repository.TransRepo
		userRepo  repository.UserRepository
	}
)

func NewTransService(tr repository.TransRepo, ur repository.UserRepository) *transService {
	return &transService{
		transRepo: tr,
		userRepo:  ur,
	}
}

func (s *transService) GetByCategory(cat string) ([]entity.Transaction, error) {
	return s.transRepo.GetByCategory(cat)
}

func (s *transService) GetAllCategoryForAdmin() ([]entity.Transaction, error) {
	return s.transRepo.GetAllCatForAdmin()
}

func (s *transService) InsertNewTrans(input entity.TransInput) error {
	return s.transRepo.InsertTrans(input)
}

func (s *transService) InsertBulkTrans(input []entity.TransInput) error {
	return s.transRepo.BulkInsertTrans(input)
}

func (s *transService) NewRecord(input entity.TransInput) error {
	userFrom, err := s.userRepo.GetuserId(input.FromId)
	if err != nil {
		fmt.Println("error inserting transaction, NewRecord, line 28")
		return err
	}

	userTo, err := s.userRepo.GetuserId(input.ToId)
	if err != nil {
		fmt.Println("error inserting transaction, NewRecord, line 34")
		return err
	}

	var transRecords []entity.TransInput

	if input.SASBalance != 0 {
		if userFrom.SASBalance == 0 || userFrom.SASBalance < input.SASBalance {
			return fmt.Errorf("error transaction, balance user %v, SASBalance : 0", input.FromId)
		} else {
			userFrom.SASBalance -= input.SASBalance
			userTo.SASBalance += input.SASBalance
		}
	}

	if input.ROBalance != 0 {
		if userFrom.ROBalance == 0 || userFrom.ROBalance < input.ROBalance {
			return fmt.Errorf("error transaction, balance user %v, ROBalance 0", input.FromId)
		} else {
			userFrom.ROBalance -= input.ROBalance
			userTo.ROBalance += input.ROBalance
		}
	}

	if input.MoneyBalance != 0 {
		if userFrom.Role == "user" {
			if userFrom.MoneyBalance == 0 || userFrom.MoneyBalance < input.MoneyBalance {
				return fmt.Errorf("error transaction, balance user %v, MoneyBalance 0", input.FromId)
			} else if userTo.Id == 1 && userTo.Role == "admin" {
				userFrom.MoneyBalance -= input.MoneyBalance
				input.Description = "pengiriman uang ke admin"
			} else {
				userFrom.MoneyBalance -= input.MoneyBalance

				transRecords = append(transRecords, entity.TransInput{
					FromId:       input.FromId,
					ToId:         1,
					Category:     entity.TransCategoryAdminFee,
					Description:  "biaya pengiriman saldo keuangan",
					MoneyBalance: entity.BiayaAdmin,
				})
			}
		}

		userTo.MoneyBalance += (input.MoneyBalance - entity.BiayaAdmin)

	}

	err = s.userRepo.UpdateBalance(userFrom)
	if err != nil {
		fmt.Println("error inserting transaction, NewRecord, line 67")
		return err
	}

	err = s.userRepo.UpdateBalance(userTo)
	if err != nil {
		fmt.Println("error inserting transaction, NewRecord, line 73")
		return err
	}

	transRecords = append(transRecords, input)

	err = s.transRepo.BulkInsertTrans(transRecords)
	if err != nil {
		fmt.Println("error inserting transaction, NewRecord, line 79")
		return err
	}

	return nil
}

func (s *transService) TransactionByUser(id int) ([]entity.Transaction, error) {
	return s.transRepo.GetTransactionById(id)
}

func (s *transService) NewDownline(inputUplineId int) error {
	var uplineId = inputUplineId

	// set true, if get parent_id 2, 3, and 4
	var checkUpperId = false

	for i := 0; i < 5; i++ {
		user, err := s.userRepo.GetuserId(uplineId)
		if err != nil {
			utils.DebugError(err, "get user id, NewDonwline")
			return err
		}

		if user.Id != 0 && user.Role != "admin" {
			var getMoney = 0

			if i == 0 {
				getMoney = 5000
			} else {
				getMoney = 3000
			}

			user.MoneyBalance += getMoney

			err := s.userRepo.UpdateBalance(user)
			if err != nil {
				utils.DebugError(err, "update balance, NewDonwline")
				return err
			}

			transInput := entity.TransInput{
				FromId:       1,
				ToId:         user.Id,
				MoneyBalance: getMoney,
				Category:     entity.TransCategoryGeneral,
				Description:  fmt.Sprintf("bonus penambahan downline baru untuk : %s", user.Fullname),
			}

			err = s.transRepo.InsertTrans(transInput)
			if err != nil {
				utils.DebugError(err, "inserting transaction, NewDonwline")
				return err
			}

			uplineId = user.ParentId

			if user.ParentId == 1 {
				checkUpperId = true
			}
		} else {
			break
		}
	}

	if !checkUpperId {
		err := s.AddBonusToUpper(uplineId)
		if err != nil {
			return err
		}
	}

	return nil
}

// give bonus with 15 level, for user id 2, 3 and 4
func (s *transService) AddBonusToUpper(startId int) error {
	var uplineId = startId

	for {
		user, err := s.userRepo.GetuserId(uplineId)
		if err != nil {
			utils.DebugError(err, " getUserId to find id upper => admin")
			return err
		}

		if user.ParentId == 1 {
			user.MoneyBalance += 3000

			err := s.userRepo.UpdateBalance(user)
			if err != nil {
				utils.DebugError(err, "update balance upper id => admin")
				return err
			}

			transInput := entity.TransInput{
				FromId:       1,
				ToId:         user.Id,
				MoneyBalance: 3000,
				Category:     entity.TransCategoryGeneral,
				Description:  fmt.Sprintf("bonus penambahan downline baru untuk : %s", user.Fullname),
			}

			err = s.transRepo.InsertTrans(transInput)
			if err != nil {
				utils.DebugError(err, "inserting transaction for upper id => admin")
				return err
			}

			break
		} else if user.Id == 1 || user.ParentId == 0 {
			break
		}

		uplineId = user.ParentId
	}

	return nil
}

func (s *transService) BuySASAdmin(input entity.BuySASAdminInput) error {
	admin, err := s.userRepo.GetuserId(1)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetuserId(input.UserId)
	if err != nil {
		return err
	}

	if admin.SASBalance < input.SASBalance {
		return errors.New("admin SAS insufficient balance")
	} else {
		admin.SASBalance -= input.SASBalance

		err := s.userRepo.UpdateBalance(admin)
		if err != nil {
			return err
		}

		user.SASBalance += input.SASBalance
		user.MoneyBalance -= input.MoneyBalance

		err = s.userRepo.UpdateBalance(user)
		if err != nil {
			return err
		}

		var transRecords []entity.TransInput

		// biaya admin dihapus
		// transRecords = append(transRecords, entity.TransInput{
		// 	FromId:       input.UserId,
		// 	ToId:         1,
		// 	Category:     entity.TransCategoryAdminFee,
		// 	Description:  "biaya admin pembelian SAS ke admin",
		// 	MoneyBalance: entity.BiayaAdmin,
		// })

		transRecords = append(transRecords, entity.TransInput{
			FromId:      1,
			ToId:        input.UserId,
			Category:    entity.TransCategorySAS,
			Description: fmt.Sprintf("pembelian SAS untuk user : %s", user.Fullname),
			SASBalance:  input.SASBalance,
		})

		err = s.transRepo.BulkInsertTrans(transRecords)
		if err != nil {
			return err
		}
	}
	return nil

}
func (s *transService) BuyROAdmin(input entity.BuyROAdminInput) error {
	admin, err := s.userRepo.GetuserId(1)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetuserId(input.UserId)
	if err != nil {
		return err
	}

	if admin.SASBalance < input.ROBalance {
		return errors.New("admin RO insufficient balance")
	} else {
		admin.ROBalance -= input.ROBalance

		err := s.userRepo.UpdateBalance(admin)
		if err != nil {
			return err
		}

		user.ROBalance += input.ROBalance
		user.MoneyBalance -= input.MoneyBalance

		err = s.userRepo.UpdateBalance(user)
		if err != nil {
			return err
		}

		var transRecord []entity.TransInput

		// biaya admin dihapus
		// transRecord = append(transRecord, entity.TransInput{
		// 	FromId:       input.UserId,
		// 	ToId:         1,
		// 	Category:     entity.TransCategoryAdminFee,
		// 	Description:  "biaya admin pembelian RO ke admin",
		// 	MoneyBalance: entity.BiayaAdmin,
		// })

		transRecord = append(transRecord, entity.TransInput{
			FromId:      1,
			ToId:        input.UserId,
			Category:    entity.TransCategoryRO,
			Description: fmt.Sprintf("pembelian RO untuk user : %s", user.Fullname),
			ROBalance:   input.ROBalance,
		})

		err = s.transRepo.BulkInsertTrans(transRecord)
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *transService) AddBalanceAdmin(input entity.AddBalanceInput) error {
	admin, err := s.userRepo.GetuserId(1)
	if err != nil {
		return err
	}

	var transRecord []entity.TransInput

	if input.ROBalance != 0 {
		transRecord = append(transRecord, entity.TransInput{
			FromId:      1,
			ToId:        1,
			Category:    entity.TransCategoryRO,
			Description: fmt.Sprintf("tambah saldo RO admin %d unit", input.ROBalance),
			ROBalance:   input.ROBalance,
		})

		admin.ROBalance += input.ROBalance
	}

	if input.SASBalance != 0 {
		transRecord = append(transRecord, entity.TransInput{
			FromId:      1,
			ToId:        1,
			Category:    entity.TransCategorySAS,
			Description: fmt.Sprintf("tambah saldo SAS admin %d unit", input.SASBalance),
			SASBalance:  input.SASBalance,
		})

		admin.SASBalance += input.SASBalance
	}

	err = s.userRepo.UpdateBalance(admin)
	if err != nil {
		return err
	}

	err = s.transRepo.BulkInsertTrans(transRecord)
	if err != nil {
		return err
	}

	return nil
}
