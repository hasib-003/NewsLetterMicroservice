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
func (h *NewsServiceHandler) GetSubscribedTopics(ctx context.Context, req *subscription.GetTopicRequest) (*subscription.GetTopicResponse, error) {
	log.Printf("Fetching topics for user ID: %v", req.GetUserId())
	topics, err := h.NewsService.GetSubscribedTopics(uint(req.GetUserId()))
	if err != nil {
		return nil, err
	}
	return &subscription.GetTopicResponse{Topics: topics}, nil
}

func RegisterNewsServiceServer(srv *grpc.Server, handler *NewsServiceHandler) {
	subscription.RegisterNewsServiceServer(srv, handler)
	log.Println("register news service success")
}
