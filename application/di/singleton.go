package di

import (
	"darbelis.eu/stabas/api"
	"darbelis.eu/stabas/util"
	"fmt"
)

var tasksRepository = NewTaskRepository()
var participantsRepository = NewParticipantsRepository()
var timeProvider = util.SimpleTimeProvider{}

var AuthenticationManager = api.NewAuthenticationManager()
var TaskControllerInstance = api.NewTaskController(tasksRepository, participantsRepository, timeProvider, AuthenticationManager)
var ParticipantsControllerInstance = api.NewParticipantController(participantsRepository, AuthenticationManager)
var AuthenticationControllerInstance = initAuthenticationController()

func initAuthenticationController() *api.AuthenticationController {
	adminPassword := AuthenticationManager.GenerateAdminPassword()
	fmt.Println("Admin password:", adminPassword)
	controller := api.NewAuthenticationController(AuthenticationManager)
	return controller
}
