package api

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

const MAX_ADMIN_TOKENS = 3

type AuthenticationController struct {
	adminPassword string
	adminTokens   []string
}

func NewAuthenticationController() *AuthenticationController {
	return &AuthenticationController{
		adminPassword: "",
		adminTokens:   []string{},
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

func (controller *AuthenticationController) GenerateAdminPassword() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rng.Intn(len(charset))]
	}

	controller.adminPassword = string(password)
	return controller.adminPassword
}

func (controller *AuthenticationController) Login(c *gin.Context) {
	var loginRequest LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}

	// Admin login
	if loginRequest.Username == "admin" {
		if loginRequest.Password != controller.adminPassword {
			c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
			return
		}

		// Check if max tokens limit reached
		if len(controller.adminTokens) >= MAX_ADMIN_TOKENS {
			c.JSON(http.StatusUnauthorized, map[string]string{"error": "Maximum number of admin sessions reached"})
			return
		}

		// Add token to adminTokens
		controller.adminTokens = append(controller.adminTokens, loginRequest.Token)

		c.JSON(http.StatusOK, LoginResponse{Token: loginRequest.Token})
		return
	}

	// TODO: Implement authentication for other users
	c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication not implemented for non-admin users"})
}
