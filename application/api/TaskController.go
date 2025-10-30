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

func (controller *TaskController) AddTask(c *gin.Context) {
	buf, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, "error reading buffer "+err.Error())
		return
	}

	task := entities.NewTask()

	err = json.API.Unmarshal(buf, &task)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error parsing json"+err.Error())
		return
	}

	controller.tasksRepository.AddTask(task)

	c.JSON(http.StatusOK, task.Id)
}

func (controller *TaskController) ChangeStatus(c *gin.Context) {
	idStr := c.Param("id")
	statusStr := c.Query("status")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusOK, "Id must be numeric "+err.Error())
	}

	status, err := strconv.Atoi(statusStr)

	if err != nil {
		c.JSON(http.StatusOK, "Status must be numeric "+err.Error())
	}

	err = entities.ValidateStatus(status)
	if err != nil {
		c.JSON(http.StatusOK, "Wrong status "+err.Error())
	}

	task := controller.tasksRepository.FindById(id)

	// TODO validate status transition

	task.Status = status

	c.JSON(http.StatusOK, "Changed status of task "+idStr+" to "+statusStr)
}

var TaskControllerInstance = TaskController{tasksRepository: dto.NewTasksRepository()}
