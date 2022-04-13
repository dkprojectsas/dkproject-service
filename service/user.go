package service

import (
	"dk-project-service/auth"
	"dk-project-service/entity"
	"dk-project-service/repository"
	"dk-project-service/utils"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type UserService interface {
	GetUserId(input int) (entity.UserDetail, error)

	GetAllUsers() ([]entity.User, error)
	GetAllUsersView(id int) ([]entity.UserView, error)

	Login(login entity.UserLogin) (entity.UserLoginResponse, error)
	Register(userAddId int, reg entity.UserRegister) error

	GetUserDownline(input string) ([]entity.User, error)

	UpdateUserById(userId int, input entity.UserUpdateInput) error
}

type userService struct {
	userRepo    repository.UserRepository
	authService auth.AuthService

	transService TransService
}

func NewUserService(userRepo repository.UserRepository, as auth.AuthService, ts TransService) *userService {
	return &userService{
		userRepo:     userRepo,
		authService:  as,
		transService: ts,
	}
}

func (s *userService) GetUserId(input int) (entity.UserDetail, error) {
	u, err := s.userRepo.GetuserId(input)

	if err != nil {
		return entity.UserDetail{}, err
	}

	userDetail := entity.UserDetail{
		Id:           u.Id,
		IdGenerate:   u.IdGenerate,
		Role:         u.Role,
		Fullname:     u.Fullname,
		PhoneNumber:  u.PhoneNumber,
		Username:     u.Username,
		ParentId:     u.ParentId,
		Position:     u.Position,
		SASBalance:   u.SASBalance,
		ROBalance:    u.ROBalance,
		MoneyBalance: u.MoneyBalance,
	}

	return userDetail, nil
}

func (s *userService) Login(login entity.UserLogin) (entity.UserLoginResponse, error) {
	var loginRes entity.UserLoginResponse

	user, err := s.userRepo.CheckUserLogin(login.Username, login.Password)

	if err != nil {
		return loginRes, err
	}

	if user.Id == 0 && user.Fullname == "" && user.Username == "" && user.PhoneNumber == "" {
		return loginRes, errors.New("error invalid data, user not registered")
	}

	userLoginToken, err := s.authService.GenerateToken(user.Id, user.Role)
	if err != nil {
		return loginRes, err
	}

	loginRes.Id = user.Id
	loginRes.Role = user.Role
	loginRes.Fullname = user.Fullname
	loginRes.PhoneNumber = user.PhoneNumber
	loginRes.ParentId = user.ParentId
	loginRes.Token = userLoginToken

	return loginRes, err
}

func (s *userService) Register(userAddId int, reg entity.UserRegister) error {
	generateId := uuid.New().String()

	parentCheck, err := s.userRepo.CheckUserId(reg.ParentId)
	if err != nil {
		return err
	}

	// var (
	// 	checkLeft, checkCenter, checkRight bool
	// )

	for i := 0; i < len(parentCheck); i++ {
		if parentCheck[i].Position == reg.Position {
			return fmt.Errorf("error downline position filled (%s)", reg.Position)
		}

		// if parentCheck[i].Position == "right" {
		// 	checkRight = true
		// }

		// if parentCheck[i].Position == "left" {
		// 	checkLeft = true
		// }

		// if parentCheck[i].Position == "center" {
		// 	checkCenter = true
		// }
	}

	// if checkCenter && checkLeft && checkRight {
	// 	return errors.New("error downline is fulfilled (left, center, right)")
	// }

	//transaction for user, using SAS balance
	parentUser, err := s.userRepo.GetuserId(userAddId)
	if err != nil {
		return err
	}

	var transRecords []entity.TransInput

	if parentUser.SASBalance < 1 {
		return errors.New("unsufficient sas balance (balance tidak cukup)")
	} else {
		parentUser.SASBalance -= 1
	}

	err = s.userRepo.UpdateBalance(parentUser)
	if err != nil {
		return err
	}

	transRecords = append(transRecords, entity.TransInput{
		FromId:      parentUser.Id,
		ToId:        1,
		SASBalance:  1,
		Category:    entity.TransCategoryGeneral,
		Description: fmt.Sprintf("pendaftaran user: %s", reg.Fullname),
	})

	var newUser entity.User

	newUser.Role = "user"
	newUser.Fullname = reg.Fullname
	newUser.PhoneNumber = reg.PhoneNumber
	newUser.ParentId = reg.ParentId
	newUser.Position = reg.Position
	newUser.Username = "DK"
	newUser.Password = utils.RandPass(6)
	newUser.IdGenerate = generateId

	createdUser, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		return fmt.Errorf("error inserting user: %s, error %s", newUser.Fullname, err.Error())
	}

	splitName := strings.Split(createdUser.Fullname, " ")

	createdUser.Username = fmt.Sprintf("DK-%v-%s", createdUser.Id, splitName[0])

	err = s.userRepo.UpdateUsername(createdUser)
	if err != nil {
		return err
	}

	err = s.transService.NewDownline(reg.ParentId)
	if err != nil {
		return err
	}

	// create send WA (concurrent)
	cbResp, err := s.userRepo.SendWANotification(createdUser)
	if err != nil {
		return err
	}

	if cbResp.MessageID != 0 && cbResp.Status != "" && cbResp.Cost != 0 {
		transRecords = append(transRecords, entity.TransInput{
			FromId:       1,
			ToId:         1,
			MoneyBalance: cbResp.Cost,
			Category:     entity.TransCategoryGeneral,
			Description:  fmt.Sprintf("notifikasi whatsapp ke user: %s", reg.Fullname),
		})
	}

	err = s.transService.InsertBulkTrans(transRecords)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) UpdateUserById(userId int, input entity.UserUpdateInput) error {
	user, err := s.userRepo.GetuserId(userId)
	if err != nil {
		return err
	}

	if user.Id == 0 && user.Fullname == "" && user.PhoneNumber == "" {
		return errors.New("user not found")
	}

	var change int = 0

	if user.Fullname != input.Fullname {
		user.Fullname = input.Fullname
		user.Username = fmt.Sprintf("DK-%d-%s", user.Id, strings.Split(input.Fullname, " ")[0])
		change++
	}

	if user.PhoneNumber != input.PhoneNumber {
		log.Println("masuk edit phone number dan password")
		user.PhoneNumber = input.PhoneNumber
		change++
	}

	if change > 0 {
		err = s.userRepo.UpdateUserById(user)
		if err != nil {
			return err
		}

		cbResp, err := s.userRepo.SendWANotification(user)
		if err != nil {
			return err
		}

		if cbResp.MessageID != 0 && cbResp.Status != "" && cbResp.Cost != 0 {
			err = s.transService.InsertNewTrans(entity.TransInput{
				FromId:       1,
				ToId:         1,
				MoneyBalance: cbResp.Cost,
				Category:     entity.TransCategoryGeneral,
				Description:  fmt.Sprintf("notifikasi whatsapp ke user: %s", input.Fullname),
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *userService) GetAllUsers() ([]entity.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *userService) GetAllUsersView(id int) ([]entity.UserView, error) {
	idStr := strconv.Itoa(id)
	return s.userRepo.GetUserViews(idStr)
}

func (s *userService) GetUserDownline(input string) ([]entity.User, error) {
	return s.userRepo.GetUsersByParentId(input)
}
