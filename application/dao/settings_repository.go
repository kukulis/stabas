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

// TODO
func (rep *ParticipantsRepository) UpdateSetting(participant *entities.Settings) error {
	// existingSetting, err := rep.FindSetting(participant.Id)

	// if err != nil {
	// 	return err
	// }

	// existingSetting.Name = participant.Name

	return nil
}
