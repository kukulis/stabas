package di

import (
	"darbelis.eu/stabas/api"
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/util"
	"fmt"
)

var AuthenticationManager *api.AuthenticationManager = nil
var TaskControllerInstance *api.TaskController = nil
var ParticipantsControllerInstance *api.ParticipantController = nil
var AuthenticationControllerInstance *api.AuthenticationController = nil
var SettingsControllerInstance *api.SettingsController = nil

func initAuthenticationController() *api.AuthenticationController {
	adminPassword := AuthenticationManager.GenerateAdminPassword()
	fmt.Println("Admin password:", adminPassword)
	controller := api.NewAuthenticationController(AuthenticationManager)
	return controller
}

func InitializeSingletons(environment string) {
	tasksRepository := NewTaskRepository(environment)
	participantsRepository := NewParticipantsRepository(environment)
	timeProvider := util.SimpleTimeProvider{}

	AuthenticationManager = api.NewAuthenticationManager(participantsRepository)
	TaskControllerInstance = api.NewTaskController(tasksRepository, participantsRepository, timeProvider, AuthenticationManager)
	ParticipantsControllerInstance = api.NewParticipantController(participantsRepository, AuthenticationManager)
	AuthenticationControllerInstance = initAuthenticationController()
	SettingsControllerInstance = api.NewSettingsController(&dao.SettingsRepository{})
}
