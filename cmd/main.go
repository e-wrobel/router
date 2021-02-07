package main

import (
	"github.com/e-wrobel/router/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	service := gin.Default()
	service.Any("/*path", handlers.HandleAnyRoute)

	err := service.Run("localhost:8080")
	if err != nil {
		panic("Unable to listen!")
	}
}
