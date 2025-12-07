package api

import (
	"darbelis.eu/stabas/entities"
	"slices"
)

func SortTasks(tasks []*entities.Task, filter TasksFilter) {
	var sortFunction = compareByNothing

	if filter.SortByTime {
		sortFunction = compareByTime
	}

	if filter.SortByStatusTime {
		sortFunction = compareByStatusTime
	}

	// Other sorting conditions

	slices.SortFunc(tasks, sortFunction)
}

func compareByNothing(a, b *entities.Task) int {
	return 0
}

func compareByTime(a, b *entities.Task) int {
	return a.CreatedAt.Compare(*b.CreatedAt)
}

func compareByStatusTime(a, b *entities.Task) int {
	return a.GetStatusTime().Compare(*b.GetStatusTime())
}
func CompareTasksById(a, b *entities.Task) int {
	if a.Id == b.Id {
		return 0
	}
	if a.Id > b.Id {
		return 1
	}
	return -1
}
