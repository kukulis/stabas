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

func NewTasksRepository(initialTasks []*entities.Task, maxId int) *TasksRepository {
	return &TasksRepository{tasks: initialTasks, maxId: maxId}
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
	if task.TaskGroup == 0 {
		task.TaskGroup = repo.maxId
	}

	repo.tasks = append(repo.tasks, task)

	return task.Id
}

func (repo *TasksRepository) UpdateTask(task *entities.Task) (*entities.Task, error) {
	t, err := repo.FindById(task.Id)

	if err != nil {
		return nil, err
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
	t.Version = task.Version

	if task.TaskGroup != 0 {
		t.TaskGroup = task.TaskGroup
	}

	return t, nil
}

func (repo *TasksRepository) DeleteTask(id int) error {
	t, err := repo.FindById(id)
	if err != nil {
		return err
	}
	t.Deleted = true

	return nil
}

func (repo *TasksRepository) UpdateTaskWithValidation(task *entities.Task) (*entities.Task, error) {
	// find existing task
	existingTask, err := repo.FindById(task.Id)

	if err != nil {
		return nil, err
	}

	if task.Version != existingTask.Version+1 {
		return existingTask, errors.New("Wrong task version")
	}

	return repo.UpdateTask(task)
}

func (repo *TasksRepository) GetCountWithSameGroup(groupId int) int {
	return len(util.ArrayFilter(repo.tasks, func(t *entities.Task) bool { return t.TaskGroup == groupId }))
}
