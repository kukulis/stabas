package dao

import (
	"darbelis.eu/stabas/entities"
	"errors"
)

type SettingsRepository struct {
	settings []*entities.Settings
}

func (rep *SettingsRepository) GetSettings() *entities.Settings {
	s, err := rep.findById(0)

	if err != nil {
		rep.AddDefaultSettings()
		return entities.NewSettings()

	}
	return s
}

func (rep *SettingsRepository) findById(id int) (*entities.Settings, error) {
	for _, setting := range rep.settings {
		if setting.Id == id {
			return setting, nil
		}
	}
	return nil, errors.New("setting not found")
}

func (rep *SettingsRepository) CheckIfExists() error {
	for _, setting := range rep.settings {
		if setting.Id == 0 {
			return nil
		}
	}
	return errors.New("setting not found")
}

func (rep *SettingsRepository) UpdateSetting(setting *entities.Settings) (*entities.Settings, error) {
	// duplicate
	s, err := rep.findById(setting.Id)

	if err != nil {
		return nil, err
	}

	*s = *setting
	//s.NewStatusDelay = setting.NewStatusDelay
	//s.NewStatusDelaySevere = setting.NewStatusDelaySevere
	//s.SentStatusDelay = setting.SentStatusDelay
	//s.SentStatusDelaySevere = setting.SentStatusDelaySevere
	//s.ReceivedStatusDelay = setting.ReceivedStatusDelay
	//s.ReceivedStatusDelaySevere = setting.ReceivedStatusDelaySevere
	//s.ExecutingStatusDelay = setting.ExecutingStatusDelay
	//s.ExecutingStatusDelaySevere = setting.ExecutingStatusDelaySevere
	//s.FinishedStatusDelay = setting.FinishedStatusDelay
	//s.FinishedStatusDelaySevere = setting.FinishedStatusDelaySevere

	return s, nil
}

func (rep *SettingsRepository) AddSettings(setting *entities.Settings) int {
	setting.Id = 0

	rep.settings = append(rep.settings, setting)

	return setting.Id
}

func (rep *SettingsRepository) AddDefaultSettings() int {
	setting := entities.NewSettings()

	rep.settings = append(rep.settings, setting)

	return setting.Id //?
}
