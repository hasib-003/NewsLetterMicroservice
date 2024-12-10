package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/repository"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/service"
	"github.com/hasib-003/newsLetterMicroservice/user-service/proto/email"
	"github.com/hasib-003/newsLetterMicroservice/user-service/proto/subscription"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	"github.com/hasib-003/newsLetterMicroservice/user-service/config"
	models "github.com/hasib-003/newsLetterMicroservice/user-service/internal/model"

	"github.com/hasib-003/newsLetterMicroservice/user-service/routes"
	"google.golang.org/grpc"
	"log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	newsconn, err := grpc.NewClient("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(newsconn *grpc.ClientConn) {
		err := newsconn.Close()
		if err != nil {

		}
	}(newsconn)
	newsClient := subscription.NewNewsServiceClient(newsconn)

	emailconn, err := grpc.NewClient("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(emailconn *grpc.ClientConn) {
		err := emailconn.Close()
		if err != nil {

		}
	}(emailconn)
	emailClient := email.NewEmailServiceClient(emailconn)

	config.ConnectDB()
	err = config.DB.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	userrepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userrepo, newsClient, emailClient)
	server := gin.Default()

	server.StaticFS("/docs", http.Dir("./docs"))
	server.GET("/swagger/*any", ginSwagger.CustomWrapHandler(
		&ginSwagger.Config{
			URL: "/docs/swagger.yaml",
		}, swaggerFiles.Handler))

	routes.RegisterRoutes(server, userService)
	err = server.Run(":8080")
	if err != nil {
		return
	}
	select {}

}
