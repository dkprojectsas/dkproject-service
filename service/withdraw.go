package service

import (
	"dk-project-service/entity"
	"dk-project-service/repository"
	"fmt"
	"log"
	"strconv"
)

type (
	WdService interface {
		GetAllWdReq() ([]entity.WdReqDetail, error)
		GetWdReqWeek() ([]entity.WdReqDetail, error)

		GetAllWdReqByUserID(userId string) ([]entity.WdReqModel, error)

		WdReqRoBalance(input entity.WdReqInput) error
		WdReqMoneyBalance(input entity.WdReqInput) error

		ApproveWdReq(id string, input entity.UpdateWdReqApprove) error
	}

	wdService struct {
		wdRepo    repository.WdRepo
		userRepo  repository.UserRepository
		transRepo repository.TransRepo
	}
)

func NewWdService(wdRepo repository.WdRepo, userRepo repository.UserRepository, transRepo repository.TransRepo) *wdService {
	return &wdService{
		wdRepo:    wdRepo,
		userRepo:  userRepo,
		transRepo: transRepo,
	}
}

func (s *wdService) ApproveWdReq(id string, input entity.UpdateWdReqApprove) error {
	wdReq, err := s.wdRepo.GetWdById(id)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetuserId(wdReq.UserId)
	if err != nil {
		return err
	}

	// update the user money saldo, with total ro_money_balance ( ro money dan jaringan)
	user.MoneyBalance += wdReq.RoMoneyBalance

	err = s.userRepo.UpdateBalance(user)
	if err != nil {
		return err
	}

	newTrans := entity.TransInput{
		FromId:       1,
		ToId:         user.Id,
		Description:  fmt.Sprintf("Penarikan RO dari user: %s", user.Username),
		SASBalance:   0,
		ROBalance:    0,
		MoneyBalance: wdReq.RoMoneyBalance,
	}

	err = s.transRepo.InsertTrans(newTrans)
	if err != nil {
		return err
	}

	return s.wdRepo.ApproveWdReqById(id, input)
}

func (s *wdService) GetAllWdReq() ([]entity.WdReqDetail, error) {
	return s.wdRepo.GetAllWd()
}

func (s *wdService) GetWdReqWeek() ([]entity.WdReqDetail, error) {
	return s.wdRepo.GetAllWdInWeek()
}

func (s *wdService) GetAllWdReqByUserID(userId string) ([]entity.WdReqModel, error) {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, err
	}

	return s.wdRepo.GetAllWdByUserID(userIdInt)
}

func (s *wdService) WdReqMoneyBalance(input entity.WdReqInput) error {
	// pengecekan money dari front end,
	// kalau saldo money lebih kecil dari yang ditarik,
	// berarti munculkan error di FE

	// check ada wd req atau nggk
	recordWdReqWeek, err := s.wdRepo.GetWdReqInWeekByUserID(input.UserId)
	if err != nil {
		return err
	}

	if recordWdReqWeek.Id == 0 && recordWdReqWeek.UserId == 0 && recordWdReqWeek.BankAccId == 0 {
		var newWdReq = entity.WdReqModel{
			UserId:    input.UserId,
			BankAccId: input.BankAccId,
		}

		if input.Moneybalance != 0 {
			newWdReq.MoneyBalance = (input.Moneybalance - entity.BiayaAdmin)
		}

		err = s.wdRepo.CreateWdReq(newWdReq)
		if err != nil {
			return err
		}
	} else {
		recordWdReqWeek.MoneyBalance += (input.Moneybalance - entity.BiayaAdmin)

		err = s.wdRepo.UpdateWdReqByID(recordWdReqWeek)
		if err != nil {
			return err
		}
	}

	user, err := s.userRepo.GetuserId(input.UserId)
	if err != nil {
		return err
	}

	user.MoneyBalance -= input.Moneybalance

	err = s.userRepo.UpdateBalance(user)
	if err != nil {
		return err
	}

	err = s.transRepo.InsertTrans(entity.TransInput{
		FromId:       input.UserId,
		ToId:         1,
		Category:     entity.TransCategoryAdminFee,
		MoneyBalance: entity.BiayaAdmin,
		Description:  "biaya admin penarikan saldo keuangan",
	})

	// pengiriman saldo ke admin
	// err = s.transRepo.InsertTrans(entity.TransInput{
	// 	FromId: input.UserId,
	// 	ToId: 1,
	// 	Category: entity.TransCategoryGeneral,
	// 	MoneyBalance: input.Moneybalance,
	// 	Description: "pengajuan pencairan",
	// })
	if err != nil {
		return err
	}

	return nil
}

func (s *wdService) WdReqRoBalance(input entity.WdReqInput) error {
	// check ro balance by front end aja
	// kalau request RO melebihi saldo , dikasih error langusng di FE

	// check ada WD req atau nggk
	recordWdReqWeek, err := s.wdRepo.GetWdReqInWeekByUserID(input.UserId)
	if err != nil {
		return err
	}

	// RUMUS :
	// data: (a) ro req, (b) ro akumulasi history, (c) ro dari upline / downline
	// si upline cair 2
	// si down 1
	// a = 1, b = 0, c = 2
	// a check c - b, jika a < c (masukkan a), jika c <= a (masukkan c), jika b >= c (tolak)
	// 1 check 2 - 0, a < b = masukkan a

	// check kebawah
	// si u4 cair 2
	// a = 1, b = 0, c = 2
	// a check c - b, a = 1, c = 2, a < c kita dapat 3000

	// si u4 cair 2
	// a = 1, b = 0, c = 1
	// a check c - b, a = 1, c = 1, c <= a kita dapat 3000

	transRecord, err := s.RoBonusNetworkUpline(input.UserId, recordWdReqWeek.RoBalance, input.RoBalance)
	if err != nil {
		return err
	}

	getBonus, err := s.GetTotalBonusNetworkDL([]int{input.UserId}, recordWdReqWeek.RoBalance, input.RoBalance)
	if err != nil {
		return err
	}

	// update user repo
	user, err := s.userRepo.GetuserId(input.UserId)
	if err != nil {
		return err
	}

	user.ROBalance -= input.RoBalance
	err = s.userRepo.UpdateBalance(user)
	if err != nil {
		return err
	}

	// total RO dan bonus
	// biaya admin dihapus
	totalROMoney := input.RoBalance * entity.BonusUser //- entity.BiayaAdmin

	// Kalau nggk ada create baru dengan logic yang ada
	if recordWdReqWeek.Id == 0 && recordWdReqWeek.UserId == 0 && recordWdReqWeek.BankAccId == 0 {
		var newWdReq = entity.WdReqModel{
			UserId:         input.UserId,
			BankAccId:      input.BankAccId,
			RoBalance:      input.RoBalance,
			RoMoneyBalance: totalROMoney + getBonus,
			Approved:       false,
		}

		err = s.wdRepo.CreateWdReq(newWdReq)
		if err != nil {
			return err
		}
	} else {
		recordWdReqWeek.RoBalance += input.RoBalance
		recordWdReqWeek.RoMoneyBalance += totalROMoney + getBonus

		err = s.wdRepo.UpdateWdReqByID(recordWdReqWeek)
		if err != nil {
			return err
		}
	}

	// 3 transaction
	// biaya admin dihapus
	// transRecord = append(transRecord, entity.TransInput{
	// 	FromId:       input.UserId,
	// 	ToId:         1,
	// 	Category:     entity.TransCategoryAdminFee,
	// 	MoneyBalance: entity.BiayaAdmin,
	// 	Description:  "biaya admin penarikan RO",
	// })

	// ro ke admin
	transRecord = append(transRecord, entity.TransInput{
		FromId:      input.UserId,
		ToId:        1,
		Category:    entity.TransCategoryGeneral,
		ROBalance:   input.RoBalance,
		Description: "kirim saldo RO untuk penarikan",
	})

	// admin kasih bonus ro money
	transRecord = append(transRecord, entity.TransInput{
		FromId:         1,
		ToId:           input.UserId,
		Category:       entity.TransCategoryGeneral,
		ROMoneyBalance: totalROMoney,
		Description:    "konversi saldo RO menjadi keuangan",
	})

	// dapat bonus jaringan
	if getBonus != 0 {
		transRecord = append(transRecord, entity.TransInput{
			FromId:         1,
			ToId:           input.UserId,
			Category:       entity.TransCategoryGeneral,
			ROMoneyBalance: getBonus,
			Description:    "Bonus jaringan match penarikan RO dari downline",
		})
	}

	err = s.transRepo.BulkInsertTrans(transRecord)
	if err != nil {
		log.Println("error inserting all transaction, wd RO user")
		return err
	}

	return nil
}

// for update ro bonus jaringan upline
func (s *wdService) RoBonusNetworkUpline(baseParentId int, inputRoAccumUser int, inputRoNowUser int) ([]entity.TransInput, error) {
	parentId := baseParentId

	// for transaction bulk insert
	var transRecords []entity.TransInput

	// looping ke atas, bonus untuk upline
	for {
		// perlu 2 data
		user, err := s.userRepo.GetuserId(parentId)
		if err != nil {
			log.Println("RoBonusNetworkUpline: error loop get data user parent bonus jaringan upline")
			return transRecords, err
		}

		parentWdReqWeek, err := s.wdRepo.GetWdReqInWeekByUserID(user.ParentId)
		if err != nil {
			log.Println("RoBonusNetworkUpline: error loop bonus wd ro upline")
			return transRecords, err
		}

		if parentWdReqWeek.Id != 0 && parentWdReqWeek.UserId != 0 && parentWdReqWeek.BankAccId != 0 {
			roCheck := parentWdReqWeek.RoBalance
			roAccumUser := inputRoAccumUser

			roNowUser := inputRoNowUser

			if roAccumUser < roCheck {
				// ro upline dikurangi akum
				roCheck -= roAccumUser

				var bonusJaringan int

				if roNowUser < roCheck {
					parentWdReqWeek.RoMoneyBalance += (roNowUser * entity.BonusJaringan)
					bonusJaringan += (roNowUser * entity.BonusJaringan)
				} else if roCheck <= roNowUser {
					parentWdReqWeek.RoMoneyBalance += (roCheck * entity.BonusJaringan)
					bonusJaringan += (roCheck * entity.BonusJaringan)
				}

				err = s.wdRepo.UpdateWdReqByID(parentWdReqWeek)
				if err != nil {
					log.Println("RoBonusNetworkUpline: error UPDATE sql loop update wd ro upline")
					return transRecords, err
				}

				transRecords = append(transRecords, entity.TransInput{
					FromId:         1,
					ToId:           parentWdReqWeek.UserId,
					Category:       entity.TransCategoryGeneral,
					Description:    "bonus jaringan penarikan RO downline",
					ROMoneyBalance: bonusJaringan})
			}
		}

		if user.ParentId == 1 {
			break
		} else {
			parentId = user.ParentId
		}
	}

	return transRecords, nil
}

// for get total bonus jaringan kita, dari check ke downline
func (s *wdService) GetTotalBonusNetworkDL(baseListId []int, inputRoAccumUser int, inputRoNowUser int) (int, error) {
	var bonus = 0

	var downlineCheckId []int = baseListId

	// loop untuk child user
	for {
		var allNextDlIds []int

		for _, dlId := range downlineCheckId {
			allDownlineUser, err := s.userRepo.CheckUserId(dlId)
			if err != nil {
				log.Println("GetTotalBonusNetworkDL: query all downline (left, center, right) loop")
				return 0, err
			}

			if len(allDownlineUser) < 1 {
				continue
			} else {
				var nextDlIds []int

				for _, dlUser := range allDownlineUser {
					dlWdReqWeek, err := s.wdRepo.GetWdReqInWeekByUserID(dlUser.Id)
					if err != nil {
						log.Println("GetTotalBonusNetworkDL: error loop bonus wd ro downline fo us")
						return 0, err
					}

					if dlWdReqWeek.Id != 0 && dlWdReqWeek.UserId != 0 && dlWdReqWeek.BankAccId != 0 {
						roDl := dlWdReqWeek.RoBalance
						roAccumUser := inputRoAccumUser

						roNowUser := inputRoNowUser

						if roAccumUser < roDl {
							roDl -= roAccumUser

							if roNowUser < roDl {
								bonus += (roNowUser * entity.BonusJaringan)
							} else if roDl <= roNowUser {
								bonus += (roDl * entity.BonusJaringan)
							}
						}
					}

					nextDlIds = append(nextDlIds, dlUser.Id)
				}

				allNextDlIds = append(allNextDlIds, nextDlIds...)
			}
		}

		if len(allNextDlIds) < 1 {
			break
		} else {
			downlineCheckId = allNextDlIds
		}
	}

	return bonus, nil
}
