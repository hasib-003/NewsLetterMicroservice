package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/service"
	"log"
	"net/http"
	"strconv"
)

type UserController struct {
	UserService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var userInput struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := c.ShouldBind(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := uc.UserService.CreateUser(userInput.Email, userInput.Name, userInput.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (uc *UserController) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	user, err := uc.UserService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (uc *UserController) GetAllUserEmails(c *gin.Context) {
	emails, err := uc.UserService.GetAllUserEmails()
	log.Println("getting all user emails from database ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": emails})
}

func (uc *UserController) SubscribeToTopic(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
		Topic string `json:"topic"`
	}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("handler topic:%v", request.Topic)
	err := uc.UserService.SubscribeToTopic(request.Email, request.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully subscribed to topic"})
}
func (uc *UserController) GetSubscribedTopic(c *gin.Context) {
	userIDParam := c.Param("user_id")

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	topics, err := uc.UserService.GetSubscribedTopics(uint32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": topics})

}
func (uc *UserController) GetSubscribedNews(c *gin.Context) {
	userIDParam := c.Param("user_id")
	log.Printf("handler topic:%v", userIDParam)
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	news, err := uc.UserService.GetSubscribedNews(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": news})
}

func (uc *UserController) SendEmails(c *gin.Context) {
	err := uc.UserService.SendEmailsToAllUsers()
	if err != nil {
		log.Printf("send emails to all users error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "emails sent"})
}
