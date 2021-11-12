package routes

import (
	"github.com/gin-gonic/gin"
	transaction "instasafe/controllers/transactions"
	"net/http"
)

//StartGin function
func StartGin() {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/transactions", transaction.GetAllTransaction)
		api.GET("/transactions/:id", transaction.GetTransaction)
		api.POST("/transactions", transaction.CreateTransaction)
		api.DELETE("/transactions/:id", transaction.DeleteTransactions)
	}
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	router.Run(":8000")
}
