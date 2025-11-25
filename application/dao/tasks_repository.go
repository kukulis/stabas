package dao

import (
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/util"
	"errors"
	"time"
)

type TasksRepository struct {
	tasks []*entities.Task
	maxId int
}

func NewTasksRepository() *TasksRepository {

	now := time.Now()
	now1 := time.Now().Add(time.Second)
	now2 := time.Now().Add(time.Second * 2)
	now3 := time.Now().Add(time.Second * 3)
	now4 := time.Now().Add(time.Second * 4)

	return &TasksRepository{
		tasks: []*entities.Task{
			{
				Id:        1,
				Message:   "task1",
				Result:    "result1",
				Sender:    1,
				Receivers: []int{2},
				Status:    1,
				Deleted:   false,
				Version:   1,
				TaskGroup: 1,
				CreatedAt: &now3,
			},
			{
				Id:        2,
				Message:   "task2",
				Result:    "result2",
				Sender:    1,
				Receivers: []int{3},
				Status:    1,
				Deleted:   false,
				Version:   1,
				TaskGroup: 1,
				CreatedAt: &now1,
			},

			{
				Id:        6,
				Message:   "task22",
				Result:    "result22",
				Sender:    1,
				Receivers: []int{3},
				Status:    1,
				Deleted:   false,
				Version:   1,
				TaskGroup: 1,
				CreatedAt: &now2,
			},
			{
				Id:        3,
				Message:   "task3",
				Result:    "result3",
				Sender:    2,
				Receivers: []int{1},
				Status:    1,
				Deleted:   false,
				Version:   1,
				TaskGroup: 3,
				CreatedAt: &now,
			},
			{
				Id:        4,
				Message:   "task4",
				Result:    "result4",
				Sender:    2,
				Receivers: []int{3},
				Status:    1,
				Deleted:   false,
				Version:   1,
				TaskGroup: 3,
				CreatedAt: &now4,
			},
		},
		maxId: 6,
	}
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

//func (repo *TasksRepository) UpdateTaskStatusAndVersion(id int, status int, version int) error {
//	t, err := repo.FindById(id)
//	if err != nil {
//		return err
//	}
//
//	t.Status = status
//	t.Version = version
//
//	return nil
//}
