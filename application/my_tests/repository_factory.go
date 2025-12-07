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
				Message:   "Pranešti apie padėtį",
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
				Message:   "Pranešti apie padėtį",
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
				Message:   "Pranešti apie padėtį",
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
				Message:   "Atsiųsti koordinates",
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
				Message:   "Atsiųsti koordinates",
				Result:    "result4",
				Sender:    2,
				Receivers: []int{3},
				Status:    1,
				Deleted:   false,
				Version:   1,
				TaskGroup: 3,
				CreatedAt: &now4,
			},
			// 'sent' task, used to fail splitting test
			{
				Id:        7,
				Message:   "Atsiųsti pastiprinimą",
				Result:    "result7",
				Sender:    2,
				Receivers: []int{3},
				Status:    2,
				Deleted:   false,
				Version:   1,
				TaskGroup: 7,
				CreatedAt: &now4,
			},
			// new task, used for splitting test
			{
				Id:        8,
				Message:   "Atsiųsti šaudmenų",
				Result:    "result8",
				Sender:    2,
				Receivers: []int{3},
				Status:    1,
				Deleted:   false,
				Version:   1,
				TaskGroup: 8,
				CreatedAt: &now4,
			},
		},
		8,
	)
}

func NewParticipantsRepository() *dao.ParticipantsRepository {
	return dao.NewParticipantsRepository(
		[]*entities.Participant{
			{Id: 1, Name: "HQ"},
			{Id: 2, Name: "KP1"},
			{Id: 3, Name: "KP2"},
			{Id: 4, Name: "KP3"},
			{Id: 5, Name: "SK1"},
			{Id: 6, Name: "SK2"},
			{Id: 7, Name: "SK3"},
			{Id: 8, Name: "SK4"},
			{Id: 9, Name: "MED"},
		},
	)
}
