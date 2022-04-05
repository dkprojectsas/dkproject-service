package main

import (
	"dk-project-service/routes"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// category transaksi

/*
INFORMATION

user need middleware, to register
login no need
getAll user to admin middleware

bank account no need middleware

transaction
- sas
- ro
- money
*/

func main() {
	r := gin.Default()

	r.Use(CORSMiddleware())

	routes.PingRoute(r)
	routes.UserRoute(r)
	routes.BankAccountRoute(r)

	routes.TransactionRoute(r)

	routes.WithdrawRoutes(r)

	r.Run()
}
