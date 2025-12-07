package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}

	// Admin login
	if loginRequest.Username == "admin" {
		err := controller.authManager.ValidateAdminLogin(loginRequest.Password, loginRequest.Token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{Token: loginRequest.Token})
		return
	}

	// TODO: Implement authentication for other users
	c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication not implemented for non-admin users"})
}
