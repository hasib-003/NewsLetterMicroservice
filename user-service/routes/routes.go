package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/handler"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/service"
	"github.com/hasib-003/newsLetterMicroservice/user-service/middleware"
)

func RegisterRoutes(router *gin.Engine, userService *service.UserService) {

	userHandler := handler.NewUserController(userService)

	router.POST("/register", userHandler.RegisterUser)
	router.POST("/verifyEmail", userHandler.VerifyUserEmail)
	router.GET("/getUserByEmail", userHandler.GetUserByEmail)
	router.GET("/buySubscription", middleware.TokenValidationMiddleware(), userHandler.BuySubscription)
	router.POST("/subscribeToTopic", middleware.TokenValidationMiddleware(), userHandler.SubscribeToTopic)
	router.GET("/getSubscribedTopic/:user_id", middleware.TokenValidationMiddleware(), userHandler.GetSubscribedTopic)
	router.GET("/getSubscribedNews", middleware.TokenValidationMiddleware(), userHandler.GetSubscribedNews)
	router.GET("/getAllUserEmails", middleware.TokenValidationMiddleware(), userHandler.GetAllUserEmails)
	router.GET("/sendEmails", middleware.TokenValidationMiddleware(), userHandler.SendEmails)
	router.GET("/publishNews", middleware.TokenValidationMiddleware(), userHandler.PublishNews)
	router.POST("/login", userHandler.Login)
	router.GET("/auth/google", userHandler.GoogleLogin)
	router.GET("/auth/callback", userHandler.GoogleCallback)

}
