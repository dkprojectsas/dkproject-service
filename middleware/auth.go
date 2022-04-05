package middleware

import (
	"dk-project-service/auth"
	"dk-project-service/utils"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Middleware(authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || len(authHeader) == 0 {
			c.AbortWithStatusJSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user unauthorize")))
			return
		}

		token, err := authService.ValidateToken(authHeader)

		if err != nil {
			c.AbortWithStatusJSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user unauthorize")))
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.AbortWithStatusJSON(401, utils.ErrorMessages(utils.ErrorUnauthorizeUser, errors.New("error user unauthorize")))
			return
		}

		userID := int(claim["user_id"].(float64))
		role := claim["role"].(string)

		// set data di handler
		c.Set("user_id", userID)
		c.Set("role", role)
	}
}
