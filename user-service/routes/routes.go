package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/user-service/middleware"

	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/handler"

	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/service"
)

func RegisterRoutes(router *gin.Engine, userService *service.UserService) {

	userHandler := handler.NewUserController(userService)

	router.POST("/register", userHandler.RegisterUser)
	router.GET("/getUserByEmail", userHandler.GetUserByEmail)
	router.POST("/subscribeToTopic", middleware.TokenValidationMiddleware(), userHandler.SubscribeToTopic)
	router.GET("/getSubscribedTopic/:user_id", middleware.TokenValidationMiddleware(), userHandler.GetSubscribedTopic)
	router.GET("/getSubscribedNews/:user_id", middleware.TokenValidationMiddleware(), userHandler.GetSubscribedNews)
	router.GET("/getAllUserEmails", middleware.TokenValidationMiddleware(), userHandler.GetAllUserEmails)
	router.GET("/sendEmails", userHandler.SendEmails)
	router.POST("/login", userHandler.Login)

}
