package controller

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	"applichic.com/chic_secret/service"
	"applichic.com/chic_secret/util"
	validator2 "applichic.com/chic_secret/validator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	guuid "github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

const codeErrorServer = "CODE_ERROR_SERVER"
const codeErrorVerificationTokenIsInvalid = "CODE_ERROR_VERIFICATION_TOKEN_INVALID"
const codeErrorEmailOrPasswordIncorrect = "CODE_ERROR_EMAIL_OR_PASSWORD_INCORRECT"

type UserClaim struct {
	User model.User
	jwt.StandardClaims
}

type AuthController struct {
	userService       *service.UserService
	tokenService      *service.TokenService
	loginTokenService *service.LoginTokenService
}

func NewAuthController() *AuthController {
	authController := new(AuthController)
	authController.userService = new(service.UserService)
	authController.tokenService = new(service.TokenService)
	authController.loginTokenService = new(service.LoginTokenService)
	return authController
}

// Create the access token with the service information
func createAccessToken(user model.User) (string, error) {
	var newUser = model.User{}
	newUser.ID = user.ID
	expiresAt := time.Now().Add(time.Duration(config.Conf.JwtTokenExpiration) * time.Minute)
	claims := UserClaim{
		newUser,
		jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	// Generates access accessToken and refresh accessToken
	unSignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return unSignedToken.SignedString([]byte(config.Conf.JwtSecret))
}

// AskCodeToLogin Ask to send a code to the email to login
func (a *AuthController) AskCodeToLogin(c *gin.Context) {
	askCodeForm := validator2.AskCodeForm{}
	if err := c.ShouldBindJSON(&askCodeForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(askCodeForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the user
	user, err := a.userService.FetchUserByEmail(askCodeForm.Email)

	if err != nil {
		user = model.User{Email: askCodeForm.Email}
		err = a.userService.Save(&user)

		// Check if there is not an error during the database query
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
				"code":  codeErrorServer,
			})

			return
		}
	}

	// Generate the code to send
	verificationCode := util.GenerateLoginToken()
	verificationCodeString := strconv.Itoa(verificationCode)

	// Save the login token in database
	dateExpiration := time.Now().Add(time.Duration(config.Conf.VerificationTokenExpiration) * time.Minute)
	loginToken := model.LoginToken{Token: verificationCode, UserID: user.ID, ExpireAt: dateExpiration}
	loginToken, err = a.loginTokenService.Save(loginToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send an email to the user with the code
	go util.SendEmail(askCodeForm.Email, "Code: "+verificationCodeString, "Chic Secret: Login to our services")
	c.JSONP(http.StatusOK, gin.H{})
}

// Login the service and send back the access token and the refresh token
func (a *AuthController) Login(c *gin.Context) {
	// Retrieve the body
	loginUserForm := validator2.LoginUserForm{}
	if err := c.ShouldBindJSON(&loginUserForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(loginUserForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user
	user, err := a.userService.FetchUserByEmail(loginUserForm.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorEmailOrPasswordIncorrect,
		})
		return
	}

	// Check if the code is exists and not expired
	_, err = a.loginTokenService.FetchTokenNotExpiredByUserId(user.ID, loginUserForm.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorVerificationTokenIsInvalid,
		})
		return
	}

	// Delete all the login tokens
	err = a.loginTokenService.DeleteAllForUser(user.ID)
	if err != nil {
		print(err)
	}

	// Create the tokens
	accessToken, err := createAccessToken(user)

	// Send an error if the tokens didn't sign well
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Retrieve the refresh token
	refreshToken, err := a.tokenService.FetchTokenByUserId(user.ID)
	if err != nil {
		uuid, errRefreshToken := guuid.NewUUID()
		if errRefreshToken != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Impossible to generate a refresh token",
				"code":  codeErrorServer,
			})
			return
		}

		refreshToken = model.Token{Token: uuid.String(), UserID: user.ID, IsValid: true}
		refreshToken, err = a.tokenService.Save(refreshToken)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
				"code":  codeErrorServer,
			})

			return
		}
	}

	// Send the tokens
	c.JSONP(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken.Token,
		"expiresIn":    config.Conf.JwtTokenExpiration,
	})
}

// RefreshAccessToken Refresh the access token thanks to a refresh token
func (a *AuthController) RefreshAccessToken(c *gin.Context) {
	// Retrieve the body
	refreshingTokenForm := validator2.RefreshingTokenForm{}
	if err := c.ShouldBindJSON(&refreshingTokenForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(refreshingTokenForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the service linked to the token
	user, err := a.userService.FetchUserFromRefreshToken(refreshingTokenForm.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to retrieve the user",
			"code":  codeErrorServer,
		})
		return
	}

	// Create the access token
	accessToken, err := createAccessToken(user)

	// Send an error if the tokens didn't sign well
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to generate the access token",
			"code":  codeErrorServer,
		})
		return
	}

	// Send the tokens
	c.JSONP(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"expiresIn":   config.Conf.JwtTokenExpiration,
	})
}
