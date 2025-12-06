package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthenticationController struct {
}

func NewAuthenticationController() *AuthenticationController {
	return &AuthenticationController{}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

	// TODO: Implement actual authentication logic
	// For now, just return an error
	c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication not implemented yet"})
}
