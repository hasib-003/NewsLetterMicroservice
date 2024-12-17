package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/service"
	"github.com/markbates/goth/gothic"
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
		Role     string `json:"role"`
	}
	if err := c.ShouldBind(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := uc.UserService.CreateUser(userInput.Email, userInput.Name, userInput.Password, userInput.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = ""
	user.VerificationToken = ""
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (uc *UserController) VerifyUserEmail(c *gin.Context) {
	var verifyRequest struct {
		Email             string `json:"email"`
		VerificationToken string `json:"verification_token"`
	}
	if err := c.ShouldBind(&verifyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user, err := uc.UserService.GetUserByEmail(verifyRequest.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if user.Verified {
		c.JSON(http.StatusConflict, gin.H{"error": "user is already verified"})
		return
	}
	err = uc.UserService.MarkEmailAsVerified(user, verifyRequest.VerificationToken)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invalid verification token" || err.Error() == "verification token is already verified" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": "failed to mark email as verified"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})

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

func (uc *UserController) Login(c *gin.Context) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBind(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	token, err := uc.UserService.Login(userInput.Email, userInput.Password)
	if err != nil {
		if err.Error() == "invalid email or password " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		if err.Error() == "user email is not verified" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user email not verified"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (uc *UserController) GoogleLogin(c *gin.Context) {
	req := c.Request
	res := c.Writer
	req.URL.RawQuery = "provider=google"
	gothic.BeginAuthHandler(res, req)
}

func (uc *UserController) GoogleCallback(c *gin.Context) {
	req := c.Request
	res := c.Writer
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to complete user auth"})
		return
	}
	token, err := uc.UserService.LoginwithGoogle(user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to login with google"})
	}
	c.JSON(http.StatusOK, gin.H{"successfully login with google . token": token})

}

func (uc *UserController) GetAllUserEmails(c *gin.Context) {
	emails, err := uc.UserService.GetAllUserEmails()
	log.Println("getting all user emails from DB ")
	log.Printf("Che%v", emails)
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
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id not found"})
	}
	log.Printf("userId:%v", userId)
	userIdStr, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
	}
	userID, err := strconv.ParseUint(userIdStr, 10, 32)
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
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role not found"})
	}
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "role is not admin"})
	}
	err := uc.UserService.SendEmailsToAllUsers()
	if err != nil {
		log.Printf("send emails to all users error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "emails sent"})
}

func (uc *UserController) PublishNews(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role not found"})
	}
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "role is not admin"})
	}
	userWithNews, err := uc.UserService.PublishUserWithNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": userWithNews})

}
