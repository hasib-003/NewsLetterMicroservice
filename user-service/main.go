package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/repository"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/service"
	"github.com/hasib-003/newsLetterMicroservice/user-service/proto/email"
	"github.com/hasib-003/newsLetterMicroservice/user-service/proto/payment"
	"github.com/hasib-003/newsLetterMicroservice/user-service/proto/subscription"
	"github.com/hasib-003/newsLetterMicroservice/user-service/utils"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"

	"github.com/hasib-003/newsLetterMicroservice/user-service/config"
	models "github.com/hasib-003/newsLetterMicroservice/user-service/internal/model"

	"github.com/hasib-003/newsLetterMicroservice/user-service/routes"
	"google.golang.org/grpc"
	"log"
)

func main() {
	gin.SetMode(gin.DebugMode)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	utils.SetupOauth()
	newsconn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(newsconn *grpc.ClientConn) {
		err := newsconn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}(newsconn)
	newsClient := subscription.NewNewsServiceClient(newsconn)

	emailconn, err := grpc.NewClient("localhost:50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(emailconn *grpc.ClientConn) {
		err := emailconn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}(emailconn)
	emailClient := email.NewEmailServiceClient(emailconn)
	paymentconn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(paymentconn *grpc.ClientConn) {
		err := paymentconn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}(paymentconn)
	paymentClient := payment.NewPaymentServiceClient(paymentconn)
	rabitmqConn := config.ConnectRabitmq()
	defer func(rabitmqConn *amqp.Connection) {
		err := rabitmqConn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}(rabitmqConn)

	config.ConnectDB()
	err = config.DB.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	userrepo := repository.NewUserRepository(config.DB, rabitmqConn)
	userService := service.NewUserService(userrepo, newsClient, emailClient, paymentClient)
	server := gin.Default()
	err = server.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}
	utils.StartCorn("http://localhost:8080/publishNews")
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
