package api

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/util"
	"errors"
	"github.com/gin-gonic/gin"
)

const MAX_ADMIN_TOKENS = 3

type AuthenticationManager struct {
	adminPassword         string
	adminTokens           []string
	CheckAuthorization    bool
	participantRepository dao.IParticipantsRepository
}

// NewAuthenticationManager creates and initializes a new AuthenticationManager instance
// with the provided participant repository
func NewAuthenticationManager(participantRepository dao.IParticipantsRepository) *AuthenticationManager {
	return &AuthenticationManager{
		adminPassword:         "",
		adminTokens:           []string{},
		CheckAuthorization:    false,
		participantRepository: participantRepository,
	}
}

// GenerateAdminPassword generates a random 10-character admin password
// using uppercase letters and digits, stores it, and returns it
func (manager *AuthenticationManager) GenerateAdminPassword() string {
	manager.adminPassword = util.StringGenerator(util.UPPER_CASE_LETTERS_AND_DIGITS, 10)

	return manager.adminPassword
}

// ValidateAdminLogin validates admin credentials and adds the token to active sessions
// Returns an error if credentials are invalid or maximum sessions limit is reached
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

// Authorize checks if the request is authorized by validating the auth_token header
// Returns true if authorized or if CheckAuthorization is disabled, false otherwise
// Sends a 401 JSON response if authentication fails
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
