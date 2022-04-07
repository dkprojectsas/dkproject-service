package routes

import (
	"dk-project-service/controller"
	"dk-project-service/repository"
	"dk-project-service/service"

	"github.com/gin-gonic/gin"
)

var (
	tsRepo       = repository.NewTransRepo(DB)
	tsService    = service.NewTransService(tsRepo, userRepo)
	tsController = controller.NewtransactionController(tsService)
)

func TransactionRoute(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/transaction", MainMiddleware, tsController.TransactionByUser)
		v1.GET("/transaction/:category", MainMiddleware, tsController.GetAllTransByCategory)
		v1.GET("/transaction/admin/transaction", MainMiddleware, tsController.GetAllTransactionForAdmin)
		v1.POST("/transaction/record", MainMiddleware, tsController.NewRecord)
		v1.POST("/transaction/buy_ro_admin", MainMiddleware, tsController.BuyROToAdmin)
		v1.POST("/transaction/buy_sas_admin", MainMiddleware, tsController.BuySASToAdmin)
		v1.POST("/transaction/add_balance_admin", MainMiddleware, tsController.AddBalanceAdmin)
	}
}
