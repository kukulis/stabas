package my_tests

import (
	"darbelis.eu/stabas/api"
	"darbelis.eu/stabas/entities"
	"github.com/gin-gonic/gin/codec/json"
	"reflect"
	"testing"
)

func TestGroupTasks(t *testing.T) {
	tasks := []*entities.Task{
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
		},
		{
			Id:        5,
			Message:   "independent",
			Result:    "independent result",
			Sender:    2,
			Receivers: []int{3},
			Status:    1,
			Deleted:   false,
			Version:   1,
			TaskGroup: 5,
		},
	}

	got := api.GroupTasks(tasks)

	want := []*entities.Task{
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
			Children: []*entities.Task{
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
				},
			},
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
			Children: []*entities.Task{
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
				},
			},
		},
		{
			Id:        5,
			Message:   "independent",
			Result:    "independent result",
			Sender:    2,
			Receivers: []int{3},
			Status:    1,
			Deleted:   false,
			Version:   1,
			TaskGroup: 5,
			Children:  []*entities.Task{},
		},
	}

	//if gotJson != wantJson {
	//	t.Errorf("Arrays were not equal want %v \n\n,  got %v", wantJson, gotJson)
	//}
	if !reflect.DeepEqual(want, got) {
		gotJsonB, _ := json.API.MarshalIndent(got, "", " ")
		wantJsonB, _ := json.API.MarshalIndent(want, "", " ")

		gotJson := string(gotJsonB)
		wantJson := string(wantJsonB)
		t.Errorf("Arrays were not equal WANT: \n %v \n\n,  GOT: \n %v\n", wantJson, gotJson)
	}
}
