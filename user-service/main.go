package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/user-service/config"
	models "github.com/hasib-003/newsLetterMicroservice/user-service/internal/model"
	"github.com/hasib-003/newsLetterMicroservice/user-service/routes"
)

func main() {
	config.ConnectDB()
	err := config.DB.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	server := gin.Default()
	routes.RegisterRoutes(server)
	err = server.Run(":8080")
	if err != nil {
		return
	}

}
