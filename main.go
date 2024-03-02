package main

import (
	"rinha/controllers"
	"rinha/database"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	database.CreateConnection()
	r.GET("/clientes/:id/extrato", controllers.Extrato)
	r.POST("/clientes/:id/transacoes", controllers.Transacoes)
	r.Run(":3000")
}
