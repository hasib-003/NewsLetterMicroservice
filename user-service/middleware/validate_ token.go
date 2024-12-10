package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

var mySignedKey []byte

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("MY_SIGNED_KEY")
	if key == "" {
		log.Fatal("Error loading MY_SIGNED_KEY")
	}
	mySignedKey = []byte(key)
}

func TokenValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is empty"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := validateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userId, ok := claims["user_id"].(string); ok {
				c.Set("userId", userId)
			}
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
		}
		c.Next()
	}
}
func validateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySignedKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
