package dao

import "darbelis.eu/stabas/entities"

type ITasksRepository interface {
	FindById(id int) (*entities.Task, error)
	FindAll() []*entities.Task
	AddTask(task *entities.Task) int
	UpdateTask(task *entities.Task) (*entities.Task, error)
	DeleteTask(id int) error
	UpdateTaskWithValidation(task *entities.Task) (*entities.Task, error)
	GetCountWithSameGroup(groupId int) int
}
