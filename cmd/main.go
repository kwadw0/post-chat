package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"kwadw0/post-chat/services"
)

func main() {
	router := gin.Default()
	message := services.ChatService("Hello from Main")
	fmt.Println(message)
	//router.GET("/", services.ChatService)
	router.Run(":8000")
}
