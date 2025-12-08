package api

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

const MAX_ADMIN_TOKENS = 3
const ADMIN_LOGIN = "admin"

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

// tryAdminLogin validates admin credentials and adds the token to active sessions
// Returns an error if credentials are invalid or maximum sessions limit is reached
func (manager *AuthenticationManager) tryAdminLogin(password string, token string) error {
	if password != manager.adminPassword {
		return errors.New("Invalid credentials")
	}

	if len(manager.adminTokens) >= MAX_ADMIN_TOKENS {
		return errors.New("Maximum number of admin sessions reached")
	}

	manager.adminTokens = append(manager.adminTokens, token)
	return nil
}

func (manager *AuthenticationManager) tryLogin(loginRequest LoginRequest) error {
	if loginRequest.Username == ADMIN_LOGIN {
		err := manager.tryAdminLogin(loginRequest.Password, loginRequest.Token)

		if err != nil {
			return err
		}

		return nil
	}

	participant := manager.participantRepository.FindParticipantByName(loginRequest.Username)

	if participant == nil {
		return errors.New("participant " + loginRequest.Username + " not found")
	}

	if participant.Token != "" {
		return errors.New("user already logged in")
	}

	if participant.Password != loginRequest.Password {
		return errors.New("wrong password")
	}

	participant.Token = loginRequest.Token

	err := manager.participantRepository.UpdateParticipantToken(participant.Id, loginRequest.Token)
	if err != nil {
		return err
	}

	return nil
}

func (manager *AuthenticationManager) Authorize(userName string, action string, context any) bool {
	if userName == ADMIN_LOGIN {
		return true
	}

	if action == "" {
		return true
	}

	// user is not admin at this point
	switch action {
	case "GetParticipants":
		return true
	case "GetParticipant":
		return true
	case "UpdateParticipant":
		return false
	case "AddParticipant":
		return false
	case "DeleteParticipant":
		return false
	case "RegeneratePassword":
		return false
	case "GetAllTasks":
		return true
	case "GetTasksGroups":
		return true
	case "GetTask":
		return true
	case "AddTask":
		return true
	case "UpdateTask":
		// TODO depends on context
		//task := context
		//if task.Sender == userName {
		//	return true
		//}

		return true
	case "DeleteTask":
		// TODO depends on context
		return true
	case "ChangeStatus":
		// TODO depends on context
		return true
	default:
		_ = fmt.Errorf("AuthenticationManager.authorize: No permission to authorize %s", action)
		return false
	}
}

func (manager *AuthenticationManager) Authenticate(c *gin.Context) (string, error) {
	if !manager.CheckAuthorization {
		return ADMIN_LOGIN, nil
	}

	authToken := c.GetHeader("auth_token")
	if authToken == "" {
		c.JSON(401, gin.H{"error": "Missing authentication token"})
		return "", errors.New("missing authentication token")
	}

	// Check if token exists in adminTokens
	for _, token := range manager.adminTokens {
		if token == authToken {
			return ADMIN_LOGIN, nil
		}
	}

	// TODO call specific method to the repository with a token value
	participants := manager.participantRepository.GetParticipants()
	for _, participant := range participants {
		if participant.Token == authToken {
			return participant.Name, nil
		}
	}

	c.JSON(401, gin.H{"error": "Invalid authentication token"})
	return "", errors.New("invalid authentication token")
}
