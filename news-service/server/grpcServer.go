package server

import (
	"github.com/hasib-003/newsLetterMicroservice/news-service/config"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/handler"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/repository"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/service"
	"github.com/hasib-003/newsLetterMicroservice/news-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func StartGrpcServer() {
	repo := &repository.NewsRepository{DB: config.DB}
	newsService := service.NewNewsService(repo)
	h := &handler.NewsServiceHandler{
		NewsService: newsService,
	}
	grpcServer := grpc.NewServer()
	subscription.RegisterNewsServiceServer(grpcServer, h)
	listener, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalln(err)
	}
	reflection.Register(grpcServer)
	log.Println("grpc server listening on :5001")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
