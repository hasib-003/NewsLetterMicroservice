package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/news-service/config"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/model"
	"github.com/hasib-003/newsLetterMicroservice/news-service/routes"
	server2 "github.com/hasib-003/newsLetterMicroservice/news-service/server"
	"log"
)

func main() {
	config.ConnectDB()
	err := config.DB.AutoMigrate(&model.News{}, &model.Topic{}, &model.Subscription{})
	if err != nil {
		panic(err)
	}
	log.Println("Starting gRPC server on :5001")
	go server2.StartGrpcServer()
	go func() {
		server := gin.Default()
		routes.RegisterNewsRoutes(server)
		err = server.Run(":8081")
		log.Println("Starting  server on :8081")
		if err != nil {
			panic(err)
		}
	}()
	select {}

}
