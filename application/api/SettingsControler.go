package api

import (
	"net/http"

	"darbelis.eu/stabas/dao"
	"github.com/gin-gonic/gin"
)

type SettingsController struct {
	settingsRepository *dao.SettingsRepository
}

func NewSettingsController(settingsRepository *dao.SettingsRepository) *SettingsController {
	return &SettingsController{settingsRepository: settingsRepository}
}

func (controller *SettingsController) GetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, controller.settingsRepository.GetSettings())
}
