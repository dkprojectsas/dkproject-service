package routes

import (
	"dk-project-service/controller"

	"github.com/gin-gonic/gin"
)

func PingRoute(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/ping", controller.Ping)
	}
}
