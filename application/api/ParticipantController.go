package api

import (
	"darbelis.eu/stabas/entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ParticipantController struct {
	participants []*entities.Participant
}

func (controller *ParticipantController) GetParticipants(c *gin.Context) {
	c.JSON(http.StatusOK, controller.participants)
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
