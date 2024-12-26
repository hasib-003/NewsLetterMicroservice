package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/payment-service/handler"
	payment "github.com/hasib-003/newsLetterMicroservice/payment-service/proto"
	"github.com/hasib-003/newsLetterMicroservice/payment-service/repository"
	"github.com/hasib-003/newsLetterMicroservice/payment-service/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	stripeKey := os.Getenv("STRIPE_KEY")
	if stripeKey == "" {
		log.Fatal("Error loading stripe key")
	}

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	repo := repository.NewStripePaymentRepository(stripeKey)
	paymentService := service.NewPaymentService(repo)
	paymentHandler := handler.NewPaymentServiceHandler(paymentService)
	grpcServer := grpc.NewServer()
	payment.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	router := gin.Default()
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Fatal("Error loading webhook secret")
	}
	webhookHandler := handler.NewWebhookHandler(webhookSecret)
	router.POST("/webhook", webhookHandler.HandleWebhook)
	go func() {
		log.Printf("Listening on :8090...")
		if err := http.ListenAndServe(":8090", router); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	log.Println("Starting grpc server...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
