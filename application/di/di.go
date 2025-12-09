package di

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/my_tests"
)

func NewTaskRepository(environment string) *dao.TasksRepository {
	if environment == "dev" {
		return my_tests.NewTasksRepository()
	}

	if environment == "empty" {
		return dao.NewTasksRepository([]*entities.Task{}, 1)
	}

	panic("wrong config for tasks repository creation")
}

func NewParticipantsRepository(environment string) *dao.ParticipantsRepository {
	if environment == "dev" {
		return my_tests.NewParticipantsRepository()
	}
	if environment == "empty" {
		return dao.NewParticipantsRepository([]*entities.Participant{})
	}

	panic("wrong config for participants repository creation")
}
