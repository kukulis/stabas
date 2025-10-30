package api

import (
	"darbelis.eu/stabas/dto"
	"darbelis.eu/stabas/entities"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/codec/json"
	"net/http"
	"strconv"
)

type TaskController struct {
	tasksRepository *dto.TasksRepository
}

func (controller *TaskController) GetAllTasks(c *gin.Context) {
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

	err = controller.tasksRepository.UpdateTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "updating task" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, task.Id)
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
	}

	status, err := strconv.Atoi(statusStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Status must be numeric " + err.Error()})
	}

	err = entities.ValidateStatus(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong status " + err.Error()})
	}

	task, err := controller.tasksRepository.FindById(id)

	// TODO validate status transition

	task.Status = status

	c.JSON(http.StatusOK, map[string]string{"error": "Changed status of task " + idStr + " to " + statusStr})
}

// TaskControllerInstance singleton
var TaskControllerInstance = TaskController{tasksRepository: dto.NewTasksRepository()}
