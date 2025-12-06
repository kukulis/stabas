package di

import (
	"darbelis.eu/stabas/api"
	"darbelis.eu/stabas/util"
)

var tasksRepository = NewTaskRepository()
var participantsRepository = NewParticipantsRepository()
var timeProvider = util.SimpleTimeProvider{}

var TaskControllerInstance = api.NewTaskController(tasksRepository, participantsRepository, timeProvider)
var ParticipantsControllerInstance = api.NewParticipantController(participantsRepository)
var AuthenticationControllerInstance = api.NewAuthenticationController()
