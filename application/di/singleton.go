package di

import (
	"darbelis.eu/stabas/api"
	"darbelis.eu/stabas/util"
	"fmt"
)

var tasksRepository = NewTaskRepository()
var participantsRepository = NewParticipantsRepository()
var timeProvider = util.SimpleTimeProvider{}

var TaskControllerInstance = api.NewTaskController(tasksRepository, participantsRepository, timeProvider)
var ParticipantsControllerInstance = api.NewParticipantController(participantsRepository)
var AuthenticationControllerInstance = initAuthenticationController()

func initAuthenticationController() *api.AuthenticationController {
	controller := api.NewAuthenticationController()
	adminPassword := controller.GenerateAdminPassword()
	fmt.Println("Admin password:", adminPassword)
	return controller
}
