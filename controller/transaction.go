package controller

import (
	"dk-project-service/entity"
	"dk-project-service/service"
	"dk-project-service/utils"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	transactionController struct {
		transService service.TransService
	}
)

func NewtransactionController(ts service.TransService) *transactionController {
	return &transactionController{transService: ts}
}

func (tc *transactionController) NewRecord(c *gin.Context) {
	// check from id harus sama dengan yang login
	idLogin, ok := c.Get("user_id")

	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	var input entity.TransInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	if idLogin.(int) != input.FromId {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("unauthorize user, cannot create transaction")))
		return
	}

	// service

	err := tc.transService.NewRecord(input)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	var res = fmt.Sprintf("success craete transaction from : %d, to %d, date : %v", input.FromId, input.ToId, time.Now())

	c.JSON(201, gin.H{
		"message": res,
	})
}

// for umum dan admin_fee with from or to admin (id = 1)
func (tc *transactionController) GetAllTransactionForAdmin(c *gin.Context) {
	_, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	role, ok := c.Get("role")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	if role.(string) != "admin" {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("user login not admin")))
		return
	}

	trans, err := tc.transService.GetAllCategoryForAdmin()
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, trans)
}

func (tc *transactionController) GetAllTransByCategory(c *gin.Context) {
	_, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	cat := c.Param("category")

	trans, err := tc.transService.GetByCategory(cat)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, trans)
}

func (tc *transactionController) TransactionByUser(c *gin.Context) {
	idLogin, ok := c.Get("user_id")

	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	res, err := tc.transService.TransactionByUser(idLogin.(int))
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, res)
}

func (tc *transactionController) BuySASToAdmin(c *gin.Context) {
	_, ok := c.Get("user_id")

	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	var input entity.BuySASAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	err := tc.transService.BuySASAdmin(input)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("success buy SAS %d unit from admin to user %d", input.SASBalance, input.UserId)})
}

func (tc *transactionController) BuyROToAdmin(c *gin.Context) {
	_, ok := c.Get("user_id")

	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	var input entity.BuyROAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	err := tc.transService.BuyROAdmin(input)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("success buy RO %d unit from admin to user %d", input.ROBalance, input.UserId)})
}

func (tc *transactionController) AddBalanceAdmin(c *gin.Context) {
	_, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}
	role, ok := c.Get("role")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	if role.(string) != "admin" {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("user login not admin")))
		return
	}

	var input entity.AddBalanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	err := tc.transService.AddBalanceAdmin(input)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, gin.H{"message": "success add balance to admin"})
}
