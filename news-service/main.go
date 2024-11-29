package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/news-service/config"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/model"
	"github.com/hasib-003/newsLetterMicroservice/news-service/routes"
)

func main() {
	config.ConnectDB()
	err := config.DB.AutoMigrate(&model.News{})
	if err != nil {
		panic(err)
	}
	server := gin.Default()
	routes.RegisterNewsRoutes(server)
	err = server.Run(":8081")
	if err != nil {
		panic(err)
	}
}
