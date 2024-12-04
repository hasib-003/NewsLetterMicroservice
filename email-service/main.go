package main

import (
	"github.com/hasib-003/newsLetterMicroservice/email-service/handler"
	email "github.com/hasib-003/newsLetterMicroservice/email-service/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	grpcServer := grpc.NewServer()
	emailHandler := handler.NewEmailHandler()
	email.RegisterEmailServiceServer(grpcServer, emailHandler)

	listener, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("listening on :50050")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
