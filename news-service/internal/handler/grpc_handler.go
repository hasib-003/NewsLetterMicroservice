package handler

import (
	"context"
	"fmt"
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

	topics, err := h.NewsService.GetSubscribedTopics(uint(req.GetUserId()))
	if err != nil {
		return nil, err
	}
	return &subscription.GetTopicResponse{Topics: topics}, nil
}
func (h *NewsServiceHandler) GetSubscribedNews(ctx context.Context, req *subscription.GetSubscribedNewsRequest) (*subscription.GetSubscribedNewsResponse, error) {
	newsItems, err := h.NewsService.GetSubscribedNews(uint(req.UserId))
	if err != nil {
		return nil, fmt.Errorf("NewsService.GetSubscribedNews err:%v", err)
	}
	return &subscription.GetSubscribedNewsResponse{
		NewsItems: newsItems,
	}, nil
}

func RegisterNewsServiceServer(srv *grpc.Server, handler *NewsServiceHandler) {
	subscription.RegisterNewsServiceServer(srv, handler)
	log.Println("register news service success")
}
