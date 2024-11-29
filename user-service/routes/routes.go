package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/user-service/config"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/handler"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/repository"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/service"
)

func RegisterRoutes(router *gin.Engine) {
	db := config.GetDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserController(userService)

	router.POST("/register", userHandler.RegisterUser)
	router.GET("/getUserByEmail", userHandler.GetUserByEmail)

}
