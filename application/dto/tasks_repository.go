package dto

import "darbelis.eu/stabas/entities"

type TasksRepository struct {
	tasks []*entities.Task
	maxId int
}

func NewTasksRepository() *TasksRepository {
	return &TasksRepository{tasks: make([]*entities.Task, 0), maxId: 0}
}

func (repo *TasksRepository) FindById(id int) *entities.Task {
	// TODO

	return nil
}

func (repo *TasksRepository) FindAll() []*entities.Task {

	return repo.tasks
}

func (repo *TasksRepository) AddTask(task *entities.Task) int {
	repo.maxId++
	task.Id = repo.maxId

	repo.tasks = append(repo.tasks, task)

	return task.Id
}

func (repo *TasksRepository) UpdateTask(task *entities.Task) {
	// TODO
}

func (repo *TasksRepository) DeleteTask(task *entities.Task) bool {
	// TODO
	return false
}
