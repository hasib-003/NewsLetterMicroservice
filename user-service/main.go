package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/repository"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/service"
	subscription "github.com/hasib-003/newsLetterMicroservice/user-service/proto"

	"github.com/hasib-003/newsLetterMicroservice/user-service/config"
	models "github.com/hasib-003/newsLetterMicroservice/user-service/internal/model"

	"github.com/hasib-003/newsLetterMicroservice/user-service/routes"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.NewClient("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	newsClient := subscription.NewNewsServiceClient(conn)

	config.ConnectDB()
	err = config.DB.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	userrepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userrepo, newsClient)
	server := gin.Default()

	routes.RegisterRoutes(server, userService)
	err = server.Run(":8080")
	if err != nil {
		return
	}
	select {}

}
