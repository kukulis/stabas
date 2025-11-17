package di

import (
	"darbelis.eu/stabas/api"
	"darbelis.eu/stabas/dao"
)

var TaskControllerInstance = api.NewTaskController(dao.NewTasksRepository())
var ParticipantsControllerInstance = api.NewParticipantController(dao.NewParticipantsRepository())
