package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/repository"
	subscription "github.com/hasib-003/newsLetterMicroservice/news-service/proto"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

type NewsService struct {
	repository *repository.NewsRepository
}

func NewNewsService(repository *repository.NewsRepository) *NewsService {
	return &NewsService{
		repository: repository,
	}
}
func (s *NewsService) FetchAndStoreNews(topic string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Error loading NEWS_API_KEY")
	}
	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&apiKey=%s", topic, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiResponse map[string]interface{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return err
	}
	news, ok := apiResponse["articles"].([]interface{})
	if !ok {
		return errors.New("error fetching news")
	}
	var newsData []map[string]interface{}
	for i, data := range news {
		if i > 5 {
			break
		}
		newsData = append(newsData, data.(map[string]interface{}))
	}
	return s.repository.SaveNews(newsData, topic)
}
func (s *NewsService) SubscribeTopic(userID uint, topicName string) (string, error) {
	log.Println("Subscribing to topic:", topicName)
	topics, err := s.repository.FindTopicByName(topicName)
	if err != nil {
		return "", err
	}
	if len(topics) == 0 {
		return "", errors.New("no topic found")
	}
	for _, topic := range topics {
		err := s.repository.CreateSubscription(userID, topic.ID)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("Successfully subscribed to topic %s", topicName), nil
}
func (s *NewsService) GetSubscribedTopics(userID uint) ([]string, error) {
	return s.repository.GetSubscribedTopicsByUserID(userID)
}
func (s *NewsService) GetSubscribedNews(userID uint) ([]*subscription.NewsItem, error) {
	subscribedNews, err := s.repository.GetUserSubscribedNews(userID)
	if err != nil {
		return nil, err
	}
	var newsItem []*subscription.NewsItem
	for _, news := range subscribedNews {
		newsItem = append(newsItem, &subscription.NewsItem{
			NewsId:      uint32(news.NewsID),
			Title:       news.Title,
			Description: news.Description,
			TopicName:   news.TopicName,
		})
	}
	return newsItem, nil
}
