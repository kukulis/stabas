package api

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/codec/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type TaskController struct {
	tasksRepository        dao.ITasksRepository
	participantsRepository dao.IParticipantsRepository
	timeProvider           util.TimeProvider
	authManager            *AuthenticationManager
}

func NewTaskController(
	tasksRepository dao.ITasksRepository,
	participantsRepository dao.IParticipantsRepository,
	timeProvider util.TimeProvider,
	authManager *AuthenticationManager,
) *TaskController {
	return &TaskController{
		tasksRepository:        tasksRepository,
		participantsRepository: participantsRepository,
		timeProvider:           timeProvider,
		authManager:            authManager,
	}
}

func (controller *TaskController) GetAllTasks(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(userName, "GetAllTasks", nil) {
		return
	}
	c.JSON(http.StatusOK, controller.tasksRepository.FindAll())
}

func (controller *TaskController) GetTasksGroups(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(userName, "GetTasksGroups", nil) {
		return
	}
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
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(userName, "GetTask", nil) {
		return
	}

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
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
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

	if !controller.authManager.Authorize(userName, "AddTask", task) {
		return
	}

	// set current date to the task
	now := controller.timeProvider.ProvideTime()
	task.CreatedAt = &now

	controller.tasksRepository.AddTask(task)
	//fmt.Printf("Added task id %d\n", task.Id)

	c.JSON(http.StatusOK, task)
}

// UpdateTask updates task
func (controller *TaskController) UpdateTask(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}

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
	//fmt.Printf("We are updating task id %d\n", id)

	_ = receivedTask.SetStatusDateIfNil(time.Now())

	existingTask, err := controller.validateTaskForUpdate(receivedTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if !controller.authManager.Authorize(userName, "UpdateTask", existingTask) {
		return
	}

	// split task to several tasks if the task has many receivers and the status is NEW
	if receivedTask.Status == entities.STATUS_NEW && len(receivedTask.Receivers) > 1 {
		//fmt.Printf("Going to split task with %d amount of receivers\n", len(receivedTask.Receivers))
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

func (controller *TaskController) validateTaskForUpdate(receivedTask *entities.Task) (*entities.Task, error) {
	existingTask, err := controller.tasksRepository.FindById(receivedTask.Id)
	if err != nil {
		return nil, err
	}

	if existingTask.Status != entities.STATUS_NEW {
		if existingTask.Sender != receivedTask.Sender {
			return nil, errors.New("Can't modify sender if task is NOT new")
		}
		if !reflect.DeepEqual(existingTask.Receivers, receivedTask.Receivers) {
			return nil, errors.New("Can't modify receivers if task is NOT new")
		}
	}

	if existingTask.Status == entities.STATUS_NEW && len(receivedTask.Receivers) > 1 {
		countWithSameGroup := controller.tasksRepository.GetCountWithSameGroup(existingTask.TaskGroup)
		if countWithSameGroup > 1 {
			if existingTask.TaskGroup == existingTask.Id {
				return nil, errors.New("Can't add more receivers as the task has children tasks")
			} else {
				return nil, errors.New("Can't add more receivers as the task has parent task")
			}
		}
	}

	return existingTask, nil
}

// DeleteTask Deletes task
func (controller *TaskController) DeleteTask(c *gin.Context) {
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(userName, "DeleteTask", nil) {
		return
	}

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
	userName, err := controller.authManager.Authenticate(c)
	if err != nil {
		return
	}
	if !controller.authManager.Authorize(userName, "ChangeStatus", nil) {
		return
	}
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

	// validate that NEW task has receivers before changing status
	if task.Status == entities.STATUS_NEW {
		if len(task.Receivers) == 0 {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot change status from NEW when the task has no receivers"})
			return
		}
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
