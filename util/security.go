package util

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	"applichic.com/chic_secret/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// GenerateLoginToken Generate a login token
func GenerateLoginToken() int {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	return rand.Intn(max-min) + min
}

// GetUserFromToken Retrieve the user from the JWT token
func GetUserFromToken(c *gin.Context) (*model.User, error) {
	token, err := GetToken(c)

	if err != nil {
		return nil, err
	}

	userService := new(service.UserService)
	userClaims := token.Claims.(jwt.MapClaims)["User"].(map[string]interface{})
	user, err := userService.FetchUserById(userClaims["ID"])

	return &user, err
}

// GetToken Get token from the Authorization header
func GetToken(c *gin.Context) (*jwt.Token, error) {
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")

	if len(splitToken) == 1 {
		return nil, errors.New("no token found")
	}

	tokenString := strings.TrimSpace(splitToken[1])

	// Check if there is a token given
	if tokenString == "" {
		return nil, errors.New("no token found")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.JwtSecret), nil
	})

	// Check if the token is correct and valid
	if err != nil || token == nil || !token.Valid {
		return nil, errors.New("no token found")
	}

	return token, nil
}

// AuthenticationRequired Retrieve the token to check if the service is authenticated
func AuthenticationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check token
		_, err := GetToken(c)

		// Check if the token is valid
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		return
	}
}
