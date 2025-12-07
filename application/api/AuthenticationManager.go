package api

import (
	"darbelis.eu/stabas/dao"
	"errors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

const MAX_ADMIN_TOKENS = 3

type AuthenticationManager struct {
	adminPassword         string
	adminTokens           []string
	CheckAuthorization    bool
	participantRepository dao.IParticipantsRepository
}

func NewAuthenticationManager(participantRepository dao.IParticipantsRepository) *AuthenticationManager {
	return &AuthenticationManager{
		adminPassword:         "",
		adminTokens:           []string{},
		CheckAuthorization:    false,
		participantRepository: participantRepository,
	}
}

func (manager *AuthenticationManager) GenerateAdminPassword() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rng.Intn(len(charset))]
	}

	manager.adminPassword = string(password)
	return manager.adminPassword
}

func (manager *AuthenticationManager) ValidateAdminLogin(password string, token string) error {
	if password != manager.adminPassword {
		return errors.New("Invalid credentials")
	}

	if len(manager.adminTokens) >= MAX_ADMIN_TOKENS {
		return errors.New("Maximum number of admin sessions reached")
	}

	manager.adminTokens = append(manager.adminTokens, token)
	return nil
}

func (manager *AuthenticationManager) Authorize(c *gin.Context) bool {
	if !manager.CheckAuthorization {
		return true
	}

	authToken := c.GetHeader("auth_token")
	if authToken == "" {
		c.JSON(401, gin.H{"error": "Missing authentication token"})
		return false
	}

	// Check if token exists in adminTokens
	for _, token := range manager.adminTokens {
		if token == authToken {
			return true
		}
	}

	c.JSON(401, gin.H{"error": "Invalid authentication token"})
	return false
}
