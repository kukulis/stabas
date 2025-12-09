package api

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const MAX_ADMIN_TOKENS = 3
const ADMIN_LOGIN = "admin"

type JSONResponder interface {
	JSON(code int, obj any)
}

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

func (manager *AuthenticationManager) Authorize(responder JSONResponder, userName string, action string, context any) bool {
	authorized := manager.authorize(userName, action, context)

	if !authorized {
		responder.JSON(http.StatusForbidden, map[string]string{"error": "Not authorized to [" + action + "] of this task"})
	}

	return authorized
}

func (manager *AuthenticationManager) authorize(userName string, action string, context any) bool {
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
		return manager.checkSender(userName, context)
	case "UpdateTask":
		return manager.checkSenderOrReceiver(userName, context)
	case "DeleteTask":
		return manager.checkSender(userName, context)
	case "ChangeStatus":
		return manager.checkSenderOrReceiver(userName, context)
	case "UpdateSettings":
		return false
	case "GetSettings":
		return false

	default:
		_ = fmt.Errorf("AuthenticationManager.authorize: No permission to authorize %s", action)
		return false
	}
}

func (manager *AuthenticationManager) checkSenderOrReceiver(userName string, context any) bool {
	task, ok := context.(*entities.Task)
	if !ok || task == nil {
		return false
	}
	participant := manager.participantRepository.FindParticipantByName(userName)
	if participant == nil {
		return false
	}
	return task.HasSenderOrReceiver(participant.Id)
}

func (manager *AuthenticationManager) checkSender(userName string, context any) bool {
	task, ok := context.(*entities.Task)
	if !ok || task == nil {
		return false
	}
	participant := manager.participantRepository.FindParticipantByName(userName)
	if participant == nil {
		return false
	}
	return task.HasSender(participant.Id)
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

	participant := manager.participantRepository.FindParticipantByToken(authToken)
	if participant != nil {
		return participant.Name, nil
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication token"})

	return "", errors.New("invalid authentication token")
}

func (manager *AuthenticationManager) GetUserId(userName string) int {
	if userName == ADMIN_LOGIN {
		return 0
	}

	participant := manager.participantRepository.FindParticipantByName(userName)
	if participant != nil {
		return participant.Id
	}

	return -1
}
