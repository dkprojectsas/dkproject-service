package routes

import (
	"dk-project-service/controller"
	"dk-project-service/repository"
	"dk-project-service/service"

	"github.com/gin-gonic/gin"
)

var (
	wrRepo       = repository.NewWdRepo(DB)
	wrService    = service.NewWdService(wrRepo, userRepo, tsRepo)
	wrController = controller.NewWdController(wrService)
)

func WithdrawRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/withdraws", MainMiddleware, wrController.GetAllWithdrawReq)            //for admin
		v1.GET("/withdraws/in_week", MainMiddleware, wrController.GetWithdrawReqInWeek) //for admin
		v1.GET("/withdraws/:user_id", MainMiddleware, wrController.GetWithdrawReqByUser)
		v1.PATCH("/withdraws/:id", MainMiddleware, wrController.PatchWithdrawReq) //update to approved
		v1.POST("/withdraws", MainMiddleware, wrController.WithdrawReq)           //create new withdraw (money, RO) for users
	}
}
