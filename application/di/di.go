package di

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/my_tests"
)

// other values are "empty", "prod"
// TODO read values from an .env file
var tasksRepositoryConfig = "testing"

func NewTaskRepository() *dao.TasksRepository {
	if tasksRepositoryConfig == "testing" {
		return my_tests.NewTasksRepository()
	}

	if tasksRepositoryConfig == "empty" {
		return dao.NewTasksRepository([]*entities.Task{}, 1)
	}

	panic("wrong config for tasks repository creation")
}

func NewParticipantsRepository() *dao.ParticipantsRepository {
	if tasksRepositoryConfig == "testing" {
		return my_tests.NewParticipantsRepository()
	}
	if tasksRepositoryConfig == "empty" {
		return dao.NewParticipantsRepository([]*entities.Participant{})
	}

	panic("wrong config for participants repository creation")
}
