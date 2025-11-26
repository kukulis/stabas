package api

import (
	"net/http"

	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/codec/json"
)

type SettingsController struct {
	settingsRepository *dao.SettingsRepository
}

func NewSettingsController(settingsRepository *dao.SettingsRepository) *SettingsController {
	return &SettingsController{settingsRepository: settingsRepository}
}

// Get the one and only settings enitity.
func (controller *SettingsController) GetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, controller.settingsRepository.GetSettings())
}

// Update the one and only settings entity
func (controller *SettingsController) UpdateSettings(c *gin.Context) {

	// Read args.
	buf, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error reading buffer " + err.Error()})
		return
	}

	receivedSetting := entities.NewSettings()

	err = json.API.Unmarshal(buf, &receivedSetting)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error parsing json" + err.Error()})
		return
	}

	// 1. Check if a settings entity exists.
	err = controller.settingsRepository.CheckIfExists()

	// a) If not: create a new settings entity.
	if err != nil {
		controller.settingsRepository.AddSettings(receivedSetting)

		c.JSON(http.StatusOK, receivedSetting)
		// newSettingId := controller.settingsRepository.AddSettings(setting)
		return //return id maybe
	}

	// b) If yes: update this settings entity.
	controller.settingsRepository.UpdateSetting(receivedSetting)
	c.JSON(http.StatusOK, receivedSetting)
}
