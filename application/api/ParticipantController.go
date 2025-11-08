package api

import (
	"darbelis.eu/stabas/entities"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/codec/json"
	"net/http"
	"strconv"
)

type ParticipantController struct {
	// TODO use repository
	participants []*entities.Participant
}

func (controller *ParticipantController) GetParticipants(c *gin.Context) {
	// TODO use repository
	c.JSON(http.StatusOK, controller.participants)
}

func (controller *ParticipantController) UpdateParticipant(c *gin.Context) {
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

	// TODO use repositry instead
	// existing participant
	var existingParticipant *entities.Participant = nil
	for _, p := range controller.participants {
		if p.Id == id {
			existingParticipant = p
			break
		}
	}

	if existingParticipant == nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "participant was not found by id " + idStr})
		return
	}

	participantDto := entities.Participant{}

	err = json.API.Unmarshal(buf, &participantDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error parsing json" + err.Error()})
		return
	}

	existingParticipant.Name = participantDto.Name
	c.JSON(http.StatusOK, "OK")
}

var ParticipantControllerInstance = ParticipantController{
	participants: []*entities.Participant{
		&entities.Participant{
			Id:   1,
			Name: "Headquarters",
		},
		&entities.Participant{
			Id:   2,
			Name: "Department 1",
		},
		&entities.Participant{
			Id:   3,
			Name: "Department 2",
		},
		&entities.Participant{
			Id:   4,
			Name: "Department 3",
		},
		&entities.Participant{
			Id:   5,
			Name: "Company I",
		},
		&entities.Participant{
			Id:   6,
			Name: "Company II",
		},
		&entities.Participant{
			Id:   7,
			Name: "Company III",
		},
	},
}
