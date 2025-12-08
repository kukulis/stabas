package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AuthenticationController struct {
	authManager *AuthenticationManager
}

func NewAuthenticationController(authManager *AuthenticationManager) *AuthenticationController {
	return &AuthenticationController{
		authManager: authManager,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (controller *AuthenticationController) Login(c *gin.Context) {
	var loginRequest LoginRequest

	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}

	err = controller.authManager.tryLogin(loginRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: loginRequest.Token})
	return
}

func (controller *AuthenticationController) User(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"id":   strconv.Itoa(controller.authManager.GetUserId(userName)),
		"name": userName,
	})
}
