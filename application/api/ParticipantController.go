package api

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/codec/json"
	"net/http"
	"strconv"
)

type ParticipantController struct {
	participantsRepository dao.IParticipantsRepository
	authManager            *AuthenticationManager
}

func NewParticipantController(participantsRepository dao.IParticipantsRepository, authManager *AuthenticationManager) *ParticipantController {
	return &ParticipantController{
		participantsRepository: participantsRepository,
		authManager:            authManager,
	}
}

func (controller *ParticipantController) GetParticipants(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(c, userName, "GetParticipants", nil) {
		return
	}

	//  hide token and password from each participant if the current login is not admin
	participants := controller.participantsRepository.GetParticipants()

	if userName != ADMIN_LOGIN {
		participantsCopy := make([]*entities.Participant, 0)

		for _, participant := range participants {
			participantCopy := entities.Participant{}
			participantCopy = *participant
			participantCopy.Token = ""
			participantCopy.Password = ""
			participantsCopy = append(participantsCopy, &participantCopy)
		}

		participants = participantsCopy
	}

	c.JSON(http.StatusOK, participants)
}

func (controller *ParticipantController) GetParticipant(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(c, userName, "GetParticipant", nil) {
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Id must be numeric " + err.Error()})
		return
	}

	participant, err := controller.participantsRepository.FindParticipant(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// hide token and password, unless the user is admin
	if userName != ADMIN_LOGIN {

		participantCopy := entities.Participant{}
		participantCopy = *participant

		participant = &participantCopy

	}

	c.JSON(http.StatusOK, participant)
}

func (controller *ParticipantController) UpdateParticipant(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(c, userName, "UpdateParticipant", nil) {
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Id must be numeric " + err.Error()})
		return
	}

	buf, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error reading buffer " + err.Error()})
		return
	}

	participantDto := &entities.Participant{}

	err = json.API.Unmarshal(buf, &participantDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error parsing json" + err.Error()})
		return
	}

	participantDto.Id = id

	err = controller.participantsRepository.UpdateParticipant(participantDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error updating participant" + err.Error()})
		return
	}

	updatedParticipant, _ := controller.participantsRepository.FindParticipant(id)

	c.JSON(http.StatusOK, updatedParticipant)
}

func (controller *ParticipantController) AddParticipant(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(c, userName, "AddParticipant", nil) {
		return
	}
	buf, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error reading buffer " + err.Error()})
		return
	}

	participantDto := &entities.Participant{}

	err = json.API.Unmarshal(buf, &participantDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error parsing json" + err.Error()})
		return
	}

	participant, err := controller.participantsRepository.AddParticipant(participantDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error updating participant" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, participant)
}

func (controller *ParticipantController) DeleteParticipant(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(c, userName, "DeleteParticipant", nil) {
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Id must be numeric " + err.Error()})
		return
	}

	err = controller.participantsRepository.RemoveParticipant(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error deleting participant" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, strconv.Itoa(id))
}

// RegeneratePassword clears the participant's token and generates a new 5-character password
// Returns the updated participant as JSON
func (controller *ParticipantController) RegeneratePassword(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(c, userName, "RegeneratePassword", nil) {
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Id must be numeric " + err.Error()})
		return
	}

	participant, err := controller.participantsRepository.FindParticipant(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	newPassword := util.StringGenerator(util.UPPER_CASE_LETTERS_AND_DIGITS, 5)

	err = controller.participantsRepository.UpdateParticipantToken(id, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error clearing token " + err.Error()})
		return
	}

	err = controller.participantsRepository.UpdateParticipantPassword(id, newPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error updating password " + err.Error()})
		return
	}

	participant.Token = ""
	participant.Password = newPassword

	c.JSON(http.StatusOK, participant)
}
