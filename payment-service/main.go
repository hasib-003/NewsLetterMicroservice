package main

import (
	"github.com/hasib-003/newsLetterMicroservice/payment-service/handler"
	payment "github.com/hasib-003/newsLetterMicroservice/payment-service/proto"
	"github.com/hasib-003/newsLetterMicroservice/payment-service/repository"
	"github.com/hasib-003/newsLetterMicroservice/payment-service/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	stripeKey := os.Getenv("STRIPE_KEY")
	repo := repository.NewStripePaymentRepository(stripeKey)
	paymentService := service.NewPaymentService(repo)
	paymentHandler := handler.NewPaymentServiceHandler(paymentService)
	grpcServer := grpc.NewServer()
	payment.RegisterPaymentServiceServer(grpcServer, paymentHandler)
	log.Println("Starting grpc server...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
