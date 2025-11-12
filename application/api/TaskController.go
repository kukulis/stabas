package api

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/codec/json"
	"net/http"
	"strconv"
	"time"
)

type TaskController struct {
	tasksRepository *dao.TasksRepository
}

func (controller *TaskController) GetAllTasks(c *gin.Context) {

	// TODO order by statuses and dates
	c.JSON(http.StatusOK, controller.tasksRepository.FindAll())
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

	c.JSON(http.StatusOK, task.Id)
}

// UpdateTask updates task
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

	task := entities.NewTask()

	err = json.API.Unmarshal(buf, &task)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "error parsing json" + err.Error()})
		return
	}

	task.Id = id

	_ = task.SetStatusDateIfNil(time.Now())

	err = controller.tasksRepository.UpdateTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "updating task" + err.Error()})
		return
	}

	// return full task
	c.JSON(http.StatusOK, task)
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

	task.Status = status

	err = task.SetStatusDate(time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Changing status " + err.Error()})
	}

	c.JSON(http.StatusOK, map[string]string{"success": "Changed status of task " + idStr + " to " + statusStr})
}

// TaskControllerInstance singleton
var TaskControllerInstance = TaskController{tasksRepository: dao.NewTasksRepository()}
