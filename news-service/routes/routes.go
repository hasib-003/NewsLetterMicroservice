package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/news-service/config"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/handler"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/repository"
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/service"
)

func RegisterNewsRoutes(router *gin.Engine) {
	db := config.GetDB()
	newsRepository := repository.NewNewsRepository(db)
	newsService := service.NewNewsService(newsRepository)
	newsHandler := handler.NewNewsHandler(newsService)

	router.GET("/fetchNews", newsHandler.FetchAndStoreNews)

}
