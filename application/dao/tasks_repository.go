package dao

import (
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/util"
	"errors"
)

type TasksRepository struct {
	tasks []*entities.Task
	maxId int
}

func NewTasksRepository() *TasksRepository {
	return &TasksRepository{tasks: make([]*entities.Task, 0), maxId: 0}
}

func (repo *TasksRepository) FindById(id int) (*entities.Task, error) {
	for _, task := range repo.tasks {
		if task.Id == id {
			return task, nil
		}
	}
	return nil, errors.New("task not found")
}

func (repo *TasksRepository) FindAll() []*entities.Task {

	return util.ArrayFilter(repo.tasks, func(t *entities.Task) bool { return !t.Deleted })
}

func (repo *TasksRepository) AddTask(task *entities.Task) int {
	repo.maxId++
	task.Id = repo.maxId

	repo.tasks = append(repo.tasks, task)

	return task.Id
}

func (repo *TasksRepository) UpdateTask(task *entities.Task) error {
	t, err := repo.FindById(task.Id)

	if err != nil {
		return err
	}

	t.Message = task.Message
	t.Result = task.Result
	t.Sender = task.Sender
	t.Receivers = task.Receivers
	t.Status = task.Status
	t.CreatedAt = task.CreatedAt
	t.SentAt = task.SentAt
	t.ReceivedAt = task.ReceivedAt
	t.ExecutingAt = task.ExecutingAt
	t.FinishedAt = task.FinishedAt
	t.ClosedAt = task.ClosedAt

	return nil
}

func (repo *TasksRepository) DeleteTask(id int) error {
	t, err := repo.FindById(id)
	if err != nil {
		return err
	}
	t.Deleted = true

	return nil
}
