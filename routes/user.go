package routes

import (
	"dk-project-service/controller"
	"dk-project-service/repository"
	"dk-project-service/service"

	"github.com/gin-gonic/gin"
)

var (
	userRepo       = repository.NewUserRepository(DB)
	userService    = service.NewUserService(userRepo, authService, tsService)
	userController = controller.NewUserController(userService)
)

func UserRoute(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.POST("/users/register", MainMiddleware, userController.Register)
		v1.POST("/users/login", userController.Login)
		v1.GET("/users", MainMiddleware, userController.GetAllUsers) // for admin
		v1.GET("/users/by_user", MainMiddleware, userController.GetAllUsersForUserView)
		v1.GET("/users/self", MainMiddleware, userController.GetUserId)
		v1.GET("/users/downline/:id", MainMiddleware, userController.GetUserDownline)
		v1.GET("/users/validate_token", MainMiddleware, userController.ValidateTokenUser)
		v1.PUT("/users/:user_id", MainMiddleware, userController.UpdateUserById)
	}
}
