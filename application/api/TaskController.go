package api

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/codec/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TaskController struct {
	tasksRepository *dao.TasksRepository
}

func NewTaskController(tasksRepository *dao.TasksRepository) *TaskController {
	return &TaskController{tasksRepository: tasksRepository}
}

func (controller *TaskController) GetAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, controller.tasksRepository.FindAll())
}

func (controller *TaskController) GetTasksGroups(c *gin.Context) {
	// get available tasks from repository
	tasks := controller.tasksRepository.FindAll()
	tasksCopy := make([]*entities.Task, len(tasks))
	for i, task := range tasks {
		tasksCopy[i] = &entities.Task{}
		*tasksCopy[i] = *task
	}
	// group them
	groupedTasks := GroupTasks(tasksCopy)
	// TODO take parameters from request
	tasksFilter := TasksFilter{SortByTime: true, SortByStatus: true}

	SortTasks(groupedTasks, tasksFilter)

	c.JSON(http.StatusOK, groupedTasks)
}

func (controller *TaskController) GetTask(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Id must be numeric " + err.Error()})
		return
	}

	t, err := controller.tasksRepository.FindById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (controller *TaskController) AddTask(c *gin.Context) {
	buf, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error reading buffer " + err.Error()})
		return
	}

	task := entities.NewTask()

	err = json.API.Unmarshal(buf, &task)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error parsing json" + err.Error()})
		return
	}

	// set current date to the task
	now := time.Now()
	task.CreatedAt = &now

	controller.tasksRepository.AddTask(task)
	fmt.Printf("Added task id %d\n", task.Id)

	c.JSON(http.StatusOK, task)
}

// UpdateTask updates task

// TODO cover controller with the test
func (controller *TaskController) UpdateTask(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Id must be numeric " + err.Error()})
		return
	}

	buf, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error reading buffer " + err.Error()})
		return
	}

	receivedTask := entities.NewTask()

	err = json.API.Unmarshal(buf, &receivedTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error parsing json" + err.Error()})
		return
	}

	receivedTask.Id = id
	fmt.Printf("We are updating task id %d\n", id)

	_ = receivedTask.SetStatusDateIfNil(time.Now())

	// split task to several tasks if the task has many receivers and the status is NEW
	if receivedTask.Status == entities.STATUS_NEW && len(receivedTask.Receivers) > 1 {
		fmt.Printf("Going to split task with %d amount of receivers\n", len(receivedTask.Receivers))
		receivedTask.TaskGroup = id
		initialReceivers := receivedTask.Receivers

		receivedTask.Receivers = []int{initialReceivers[0]}

		// create multiple additional tasks with the same group
		receivedTask, _ := controller.tasksRepository.UpdateTask(receivedTask)

		for i := 1; i < len(initialReceivers); i++ {
			receiver := initialReceivers[i]

			additionalTask := &entities.Task{}
			// copy fields
			*additionalTask = *receivedTask

			// reset id
			additionalTask.Id = 0
			// assign receiver
			additionalTask.Receivers = []int{receiver}
			additionalTask.TaskGroup = receivedTask.TaskGroup

			controller.tasksRepository.AddTask(additionalTask)

			receivedTask.Children = append(receivedTask.Children, additionalTask)
		}
		existingTask, err := controller.tasksRepository.UpdateTask(receivedTask)

		if err != nil {
			// TODO try to reuse code block
			if strings.Contains(err.Error(), "version") {
				c.JSON(http.StatusConflict, existingTask)
				return
			}
			c.JSON(http.StatusBadRequest, map[string]string{"error": "updating receivedTask" + err.Error()})
			return
		}

		c.JSON(http.StatusOK, existingTask)
	} else {

		// TODO If the status is different than "NEW" do not let it have multiple receivers
		// also do not let task to have multiple receivers, if it has child or parent.

		existingTask, err := controller.tasksRepository.UpdateTaskWithValidation(receivedTask)
		if err != nil {
			// TODO try to reuse code block
			if strings.Contains(err.Error(), "version") {
				c.JSON(http.StatusConflict, existingTask)
				return
			}
			c.JSON(http.StatusBadRequest, map[string]string{"error": "updating receivedTask" + err.Error()})
			return
		}
		c.JSON(http.StatusOK, existingTask)
	}

}

// DeleteTask Deletes task
func (controller *TaskController) DeleteTask(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Id must be numeric " + err.Error()})
	}

	err = controller.tasksRepository.DeleteTask(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	c.JSON(http.StatusOK, id)
}

func (controller *TaskController) ChangeStatus(c *gin.Context) {
	idStr := c.Param("id")
	statusStr := c.Query("status")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Id must be numeric " + err.Error()})
		return
	}

	status, err := strconv.Atoi(statusStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Status must be numeric " + err.Error()})
		return
	}

	err = entities.ValidateStatus(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong status " + err.Error()})
		return
	}

	task, err := controller.tasksRepository.FindById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Task not found by " + idStr})
		return
	}

	// validate status transition
	if status != task.Status+1 {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Cant change status from " + strconv.Itoa(task.Status) + " to " + strconv.Itoa(status)})
		return
	}

	taskCopy := *task
	taskCopy.Status = status
	taskCopy.Version = task.Version + 1
	err = taskCopy.SetStatusDate(time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Changing status " + err.Error()})
		return
	}

	existingTask, err := controller.tasksRepository.UpdateTaskWithValidation(&taskCopy)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Changing status " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingTask)
}
