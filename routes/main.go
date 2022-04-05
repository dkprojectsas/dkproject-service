package routes

import (
	"dk-project-service/auth"
	"dk-project-service/config"
	"dk-project-service/middleware"
)

var (
	DB             = config.Conn()
	authService    = auth.NewAuthService()
	MainMiddleware = middleware.Middleware(authService)
)
