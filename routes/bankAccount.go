package routes

import (
	"dk-project-service/controller"
	"dk-project-service/repository"
	"dk-project-service/service"

	"github.com/gin-gonic/gin"
)

var (
	baRepo       = repository.NewBankAccountRepo(DB)
	baService    = service.NewBankAccountService(baRepo)
	baController = controller.NewBankAccountController(baService)
)

func BankAccountRoute(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/bank_account", MainMiddleware, baController.GetBankAccountUser)
		v1.PUT("/bank_account/:bank_acc_id", MainMiddleware, baController.UpdateBankRecord)
		v1.POST("/bank_account", MainMiddleware, baController.InsertNewBankRecord)
	}
}
