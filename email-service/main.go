package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/email-service/handler"
	email "github.com/hasib-003/newsLetterMicroservice/email-service/proto"
	"github.com/hasib-003/newsLetterMicroservice/email-service/repositoty"
	"github.com/hasib-003/newsLetterMicroservice/email-service/service"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	emailRepository := repositoty.NewRabbitMqRepository(conn)
	emailService := service.NewEmailService(emailRepository)
	emailHandler := handler.NewEmailHandler(emailService)

	go func() {
		listener, err := net.Listen("tcp", ":50050")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		email.RegisterEmailServiceServer(grpcServer, emailHandler)
		log.Println("listening on :50050")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	r := gin.Default()
	r.POST("/start", emailHandler.StartListening)
	err = r.Run(":8085")
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
	log.Println("server start on port 8085")
}
