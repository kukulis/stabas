package api

import "darbelis.eu/stabas/entities"

// GroupTasks Groups array of tasks bet putting each to others Children
func GroupTasks(tasks []*entities.Task) []*entities.Task {
	result := []*entities.Task{}

	tasksGroupsMap := make(map[int]*entities.Task)

	for _, task := range tasks {
		groupingTask, ok := tasksGroupsMap[task.TaskGroup]

		if !ok {
			groupingTask = task
			tasksGroupsMap[task.TaskGroup] = groupingTask
			groupingTask.Children = []*entities.Task{}
			continue
		}

		groupingTask.Children = append(groupingTask.Children, task)
	}

	for _, t := range tasksGroupsMap {
		result = append(result, t)
	}

	return result
}
