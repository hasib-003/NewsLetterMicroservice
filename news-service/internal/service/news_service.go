package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/repository"
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
	defer resp.Body.Close()
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
		return errors.New("Error fetching news")
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
