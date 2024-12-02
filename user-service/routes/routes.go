package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/handler"

	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/service"
)

func RegisterRoutes(router *gin.Engine, userService *service.UserService) {

	userHandler := handler.NewUserController(userService)

	router.POST("/register", userHandler.RegisterUser)
	router.GET("/getUserByEmail", userHandler.GetUserByEmail)
	router.POST("/subscribeToTopic", userHandler.SubscribeToTopic)

}
