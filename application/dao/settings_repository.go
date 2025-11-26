package dao

import (
	"errors"

	"darbelis.eu/stabas/entities"
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

func (repo *SettingsRepository) findById(id int) (*entities.Settings, error) {
	for _, setting := range repo.settings {
		if setting.Id == id {
			return setting, nil
		}
	}
	return nil, errors.New("setting not found")
}

func (repo *SettingsRepository) CheckIfExists() error {
	for _, setting := range repo.settings {
		if setting.Id == 0 {
			return nil
		}
	}
	return errors.New("setting not found")
}

func (repo *SettingsRepository) UpdateSetting(setting *entities.Settings) (*entities.Settings, error) {
	// duplicate
	s, err := repo.findById(setting.Id)

	if err != nil {
		return nil, err
	}

	s.NewStatusDelay = setting.NewStatusDelay
	s.NewStatusDelaySevere = setting.NewStatusDelaySevere
	s.SentStatusDelay = setting.SentStatusDelay
	s.SentStatusDelaySevere = setting.SentStatusDelaySevere
	s.ReceivedStatusDelay = setting.ReceivedStatusDelay
	s.ReceivedStatusDelaySevere = setting.ReceivedStatusDelaySevere
	s.ExecutingStatusDelay = setting.ExecutingStatusDelay
	s.ExecutingStatusDelaySevere = setting.ExecutingStatusDelaySevere
	s.FinishedStatusDelay = setting.FinishedStatusDelay
	s.FinishedStatusDelaySevere = setting.FinishedStatusDelaySevere

	return s, nil
}

func (repo *SettingsRepository) AddSettings(setting *entities.Settings) int {
	setting.Id = 0

	repo.settings = append(repo.settings, setting)

	return setting.Id
}

func (repo *SettingsRepository) AddDefaultSettings() int {
	setting := entities.NewSettings()

	repo.settings = append(repo.settings, setting)

	return setting.Id //?
}
