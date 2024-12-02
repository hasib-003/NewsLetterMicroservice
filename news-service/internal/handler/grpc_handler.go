package handler

import (
	"context"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/service"
	subscription "github.com/hasib-003/newsLetterMicroservice/news-service/proto"
	"google.golang.org/grpc"
	"log"
)

type NewsServiceHandler struct {
	subscription.UnimplementedNewsServiceServer
	NewsService *service.NewsService
}

func (h *NewsServiceHandler) SubscribeToTopic(ctx context.Context, req *subscription.SubscribeRequest) (*subscription.SubscribeResponse, error) {
	log.Printf("topic name:%v", req.TopicName)
	message, err := h.NewsService.SubscribeTopic(uint(req.GetUserId()), req.GetTopicName())

	if err != nil {
		return nil, err
	}
	return &subscription.SubscribeResponse{
		Success: true,
		Message: message,
	}, nil
}

func RegisterNewsServiceServer(srv *grpc.Server, handler *NewsServiceHandler) {
	subscription.RegisterNewsServiceServer(srv, handler)
	log.Println("register news service success")
}
