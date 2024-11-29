package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/service"
	"net/http"
)

type NewsHandler struct {
	newsService *service.NewsService
}

func NewNewsHandler(newsService *service.NewsService) *NewsHandler {
	return &NewsHandler{newsService: newsService}
}
func (handler *NewsHandler) FetchAndStoreNews(c *gin.Context) {
	topic := c.Query("topic")
	if topic == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Topic is required",
		})
	}
	err := handler.newsService.FetchAndStoreNews(topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "news Fetched And stored successfully"})

}
