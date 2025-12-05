package my_tests

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"time"
)

func NewTasksRepository() *dao.TasksRepository {
	now := time.Now()
	now1 := time.Now().Add(time.Second)
	now2 := time.Now().Add(time.Second * 2)
	now3 := time.Now().Add(time.Second * 3)
	now4 := time.Now().Add(time.Second * 4)

	return dao.NewTasksRepository(
		[]*entities.Task{
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
		6,
	)
}

func NewParticipantsRepository() *dao.ParticipantsRepository {
	return dao.NewParticipantsRepository(
		[]*entities.Participant{
			{Id: 1, Name: "Participant 1"},
			{Id: 2, Name: "Participant 2"},
			{Id: 3, Name: "Participant 3"},
		},
	)
}
