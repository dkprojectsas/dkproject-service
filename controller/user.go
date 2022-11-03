package controller

import (
	"dk-project-service/entity"
	"dk-project-service/service"
	"dk-project-service/utils"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *userController {
	return &userController{userService: userService}
}

func (uc *userController) ValidateTokenUser(c *gin.Context) {
	_, ok := c.Get("user_id")

	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	c.JSON(200, gin.H{
		"message": "token valid",
	})
}

func (uc *userController) GetUserId(c *gin.Context) {
	// search menggunakan user yang lagi login
	idLogin, ok := c.Get("user_id")

	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	uDetail, err := uc.userService.GetUserId(idLogin.(int))
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, uDetail)
}

func (uc *userController) Login(c *gin.Context) {
	var input entity.UserLogin

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	res, err := uc.userService.Login(input)

	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, res)
}

func (uc *userController) Register(c *gin.Context) {
	idLogin, ok := c.Get("user_id")

	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	var input entity.UserRegister

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	err := uc.userService.Register(idLogin.(int), input)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(201, gin.H{
		"message": "success register user",
	})
}

func (uc *userController) GetAllUsersForUserView(c *gin.Context) {
	idLogin, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	users, err := uc.userService.GetAllUsersView(idLogin.(int))
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, users)
}

func (uc *userController) GetAllUsers(c *gin.Context) {
	_, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	roleLogin, ok := c.Get("role")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	if roleLogin.(string) != "admin" {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user login not admin")))
		return
	}

	users, err := uc.userService.GetAllUsers()
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, users)
}

func (uc *userController) GetUserDownline(c *gin.Context) {
	_, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	IdParam := c.Param("id")

	if IdParam == "" {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, errors.New("parameter not valid")))
		return
	}

	users, err := uc.userService.GetUserDownline(IdParam)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, users)
}

func (uc *userController) UpdateUserById(c *gin.Context) {
	_, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	roleLogin, ok := c.Get("role")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	if roleLogin.(string) != "admin" {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user login not admin")))
		return
	}

	IdParam := c.Param("user_id")

	if IdParam == "" {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, errors.New("parameter not valid")))
		return
	}

	var input entity.UserUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	userId, err := strconv.Atoi(IdParam)
	if err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, errors.New("parameter not valid")))
		return
	}

	err = uc.userService.UpdateUserById(userId, input)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, gin.H{
		"message": "success update user",
	})

}

func (uc *userController) ForgetPassword(c *gin.Context) {
	// get parameter username, phone numbe
	var forgotPass entity.InputForgotPass

	if err := c.ShouldBindJSON(&forgotPass); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	res, err := uc.userService.ForgotPassword(forgotPass)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, res)
}

func (uc *userController) ChangePassword(c *gin.Context) {
	user, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	var input entity.InputChangePass

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	err := uc.userService.UpdatePasswordById(user.(int), input)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, gin.H{
		"message": "success update password",
	})
}
