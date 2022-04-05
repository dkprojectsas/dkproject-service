package controller

import (
	"dk-project-service/entity"
	"dk-project-service/service"
	"dk-project-service/utils"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	wdController struct {
		wdService service.WdService
	}
)

func NewWdController(wdService service.WdService) *wdController {
	return &wdController{wdService: wdService}
}

func (wc *wdController) GetAllWithdrawReq(c *gin.Context) {
	id, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	if id.(int) != 1 {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not admin")))
		return
	}

	res, err := wc.wdService.GetAllWdReq()
	if err != nil {
		if err != nil {
			c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
			return
		}
	}

	c.JSON(200, res)
}

func (wc *wdController) PatchWithdrawReq(c *gin.Context) {
	roleLogin, ok := c.Get("role")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	if roleLogin.(string) != "admin" {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("user not admin")))
		return
	}

	IdParam := c.Param("id")

	if IdParam == "" {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, errors.New("parameter not valid")))
		return
	}

	var input entity.UpdateWdReqApprove

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	err := wc.wdService.ApproveWdReq(IdParam, input)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("success patch wd_req id : %s , approver: %t", IdParam, input.Approved),
	})

}

func (wc *wdController) GetWithdrawReqInWeek(c *gin.Context) {
	id, ok := c.Get("user_id")

	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	if id.(int) != 1 {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not admin")))
		return
	}

	res, err := wc.wdService.GetWdReqWeek()
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, res)
}

func (wc *wdController) GetWithdrawReqByUser(c *gin.Context) {
	loginId, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	IdParam := c.Param("user_id")

	if IdParam == "" {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, errors.New("parameter not valid")))
		return
	}

	idParamInt, _ := strconv.Atoi(IdParam)

	if idParamInt != loginId.(int) {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not authorize access")))
		return
	}

	res, err := wc.wdService.GetAllWdReqByUserID(IdParam)
	if err != nil {
		c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
		return
	}

	c.JSON(200, res)
}

func (wc *wdController) WithdrawReq(c *gin.Context) {
	loginId, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user not login")))
		return
	}

	var input entity.WdReqInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, utils.ErrorMessages(utils.ErrorBadRequest, err))
		return
	}

	if loginId.(int) != input.UserId {
		c.JSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error not authorize for login user")))
		return
	}

	if input.Moneybalance != 0 {
		err := wc.wdService.WdReqMoneyBalance(input)
		if err != nil {
			c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
			return
		}
	}

	if input.RoBalance != 0 {
		err := wc.wdService.WdReqRoBalance(input)
		if err != nil {
			c.JSON(500, utils.ErrorMessages(utils.ErrorInternalServer, err))
			return
		}
	}

	c.JSON(201, gin.H{
		"message": fmt.Sprintf("success withraw request for user id : %v", input.UserId),
	})
}
