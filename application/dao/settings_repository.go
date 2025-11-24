package dao

import (
	"darbelis.eu/stabas/entities"
)

type SettingsRepository struct {
	settings []*entities.Settings
}

func (rep *SettingsRepository) GetSettings() []*entities.Settings {
	return rep.settings
}

// func (controller *SettingsController) UpdateSettings(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, map[string]string{"error": "Id must be numeric " + err.Error()})
// 		return
// 	}

// 	buf, err := c.GetRawData()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, map[string]string{"error": "error reading buffer " + err.Error()})
// 		return
// 	}

// 	participantDto := &entities.Participant{}

// 	err = json.API.Unmarshal(buf, &participantDto)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, map[string]string{"error": "error parsing json" + err.Error()})
// 		return
// 	}

// 	participantDto.Id = id

// 	err = controller.participantsRepository.UpdateParticipant(participantDto)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, map[string]string{"error": "error updating participant" + err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, "OK")
// }
